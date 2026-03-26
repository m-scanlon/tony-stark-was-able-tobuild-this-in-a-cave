#!/usr/bin/env python3
import argparse
from dataclasses import dataclass
from pathlib import Path
from typing import Dict, List, Optional, Tuple

import spacy
from fastcoref import FCoref
from spacy.matcher import DependencyMatcher
from spacy.tokens import Doc, Span, Token

from coref_resolve import force_offline_mode, preprocess_text, resolve_model_path

try:
    import jamspell
except ImportError:
    jamspell = None

NON_REWRITE_PRONOUNS = {
    "i", "me", "my", "myself",
    "you", "your", "yours", "yourself",
    "we", "us", "our", "ours", "ourselves",
}


@dataclass
class Projection:
    subject: str
    relation: str
    object: str
    pattern: str


@dataclass
class CorefAnnotation:
    clusters: List[List[Tuple[int, int]]]
    span_to_anchor: Dict[Tuple[int, int], str]


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    parser.add_argument("text", nargs="+")
    parser.add_argument("--repair", action="store_true")
    parser.add_argument("--repair-model", default="models/en.bin")
    parser.add_argument("--resolve-coref", action="store_true")
    parser.add_argument("--coref-model", default="FCoref")
    parser.add_argument("--device", default="cpu")
    return parser


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()

    text = " ".join(args.text)
    preprocessed = preprocess_text(text)
    repaired = preprocessed
    coref_annotation = CorefAnnotation(clusters=[], span_to_anchor={})

    if args.repair:
        repaired = run_repair(preprocessed, args.repair_model)

    if args.resolve_coref:
        coref_annotation = run_coref(repaired, args.coref_model, args.device)

    nlp = spacy.load("en_core_web_sm")
    doc = nlp(repaired)
    matcher = build_matcher(nlp)

    print("INPUT")
    print(text)
    print()

    if preprocessed != text:
        print("PREPROCESSED")
        print(preprocessed)
        print()

    if repaired != preprocessed:
        print("REPAIRED")
        print(repaired)
        print()

    if coref_annotation.clusters:
        print("COREFERENCE CLUSTERS")
        for index, cluster in enumerate(coref_annotation.clusters, start=1):
            rendered = []
            for start, end in cluster:
                rendered.append(repaired[start:end])
            print(f"  CLUSTER {index}: {rendered}")
        print()

    for sent_index, sent in enumerate(doc.sents, start=1):
        print(f"SENTENCE {sent_index}")
        print(sent.text)
        print("  TOKENS")
        for token in sent:
            print(
                f"    {token.i}: text={token.text} lemma={token.lemma_} pos={token.pos_} dep={token.dep_} head={token.head.i}:{token.head.text}"
            )

        print("  NOUN CHUNKS")
        noun_chunks = noun_chunk_map(sent)
        if noun_chunks:
            for root_i, chunk in noun_chunks.items():
                print(f"    root={root_i}:{sent.doc[root_i].text} chunk={chunk.text}")
        else:
            print("    none")

        print("  DEPENDENCY MATCHES")
        matches = matcher(sent.as_doc())
        if matches:
            sent_doc = sent.as_doc()
            for match_id, token_ids in matches:
                name = sent_doc.vocab.strings[match_id]
                tokens = [sent_doc[i].text for i in token_ids]
                print(f"    {name}: {tokens}")
        else:
            print("    none")

        print("  PROJECTED FRAGMENTS")
        projections = project_sentence(sent, coref_annotation)
        if projections:
            for projection in projections:
                print(
                    f"    {projection.subject} -> {projection.relation} -> {projection.object} [{projection.pattern}]"
                )
        else:
            print("    none")
        print()

    return 0


def run_repair(text: str, model_path_str: str) -> str:
    if jamspell is None:
        raise SystemExit("JamSpell is not installed in this environment.")

    model_path = Path(model_path_str)
    if not model_path.exists():
        raise SystemExit(f"Missing JamSpell model: {model_path}")

    corrector = jamspell.TSpellCorrector()
    if not corrector.LoadLangModel(str(model_path)):
        raise SystemExit(f"Failed to load JamSpell model: {model_path}")

    return corrector.FixFragment(text)


def run_coref(text: str, model_name: str, device: str) -> CorefAnnotation:
    force_offline_mode()
    model = FCoref(model_name_or_path=resolve_model_path(model_name), device=device)
    result = model.predict(texts=[text])[0]
    raw_clusters = result.get_clusters(as_strings=False)
    clusters: List[List[Tuple[int, int]]] = []
    span_to_anchor: Dict[Tuple[int, int], str] = {}

    for cluster in raw_clusters:
        normalized_cluster: List[Tuple[int, int]] = []
        for mention in cluster:
            if isinstance(mention, (list, tuple)) and len(mention) == 2:
                start, end = int(mention[0]), int(mention[1])
                normalized_cluster.append((start, end))
        if not normalized_cluster:
            continue

        anchor_start, anchor_end = normalized_cluster[0]
        anchor_text = text[anchor_start:anchor_end]
        clusters.append(normalized_cluster)
        for span in normalized_cluster:
            span_to_anchor[span] = anchor_text

    return CorefAnnotation(clusters=clusters, span_to_anchor=span_to_anchor)


def build_matcher(nlp) -> DependencyMatcher:
    matcher = DependencyMatcher(nlp.vocab)

    matcher.add(
        "VERB_DIRECT_OBJECT",
        [[
            {"RIGHT_ID": "verb", "RIGHT_ATTRS": {"POS": {"IN": ["VERB", "AUX"]}}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "subject", "RIGHT_ATTRS": {"DEP": {"IN": ["nsubj", "nsubjpass"]}}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "object", "RIGHT_ATTRS": {"DEP": {"IN": ["dobj", "obj"]}}},
        ]],
    )

    matcher.add(
        "VERB_PREP_OBJECT",
        [[
            {"RIGHT_ID": "verb", "RIGHT_ATTRS": {"POS": {"IN": ["VERB", "AUX"]}}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "subject", "RIGHT_ATTRS": {"DEP": {"IN": ["nsubj", "nsubjpass"]}}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "prep", "RIGHT_ATTRS": {"DEP": "prep"}},
            {"LEFT_ID": "prep", "REL_OP": ">", "RIGHT_ID": "pobj", "RIGHT_ATTRS": {"DEP": "pobj"}},
        ]],
    )

    matcher.add(
        "COPULA_ATTRIBUTE",
        [[
            {"RIGHT_ID": "copula", "RIGHT_ATTRS": {"LEMMA": "be", "POS": {"IN": ["AUX", "VERB"]}}},
            {"LEFT_ID": "copula", "REL_OP": ">", "RIGHT_ID": "subject", "RIGHT_ATTRS": {"DEP": {"IN": ["nsubj", "nsubjpass"]}}},
            {"LEFT_ID": "copula", "REL_OP": ">", "RIGHT_ID": "attribute", "RIGHT_ATTRS": {"DEP": "attr"}},
        ]],
    )

    matcher.add(
        "COPULA_PROPERTY",
        [[
            {"RIGHT_ID": "copula", "RIGHT_ATTRS": {"LEMMA": "be", "POS": {"IN": ["AUX", "VERB"]}}},
            {"LEFT_ID": "copula", "REL_OP": ">", "RIGHT_ID": "property", "RIGHT_ATTRS": {"DEP": "acomp"}},
        ]],
    )

    matcher.add(
        "VERB_CCOMP",
        [[
            {"RIGHT_ID": "verb", "RIGHT_ATTRS": {"POS": {"IN": ["VERB", "AUX"]}}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "subject", "RIGHT_ATTRS": {"DEP": {"IN": ["nsubj", "nsubjpass"]}}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "ccomp", "RIGHT_ATTRS": {"DEP": {"IN": ["ccomp", "xcomp"]}}},
            {"LEFT_ID": "ccomp", "REL_OP": ">", "RIGHT_ID": "ccomp_subject", "RIGHT_ATTRS": {"DEP": {"IN": ["nsubj", "nsubjpass"]}}},
        ]],
    )

    matcher.add(
        "VERB_XCOMP",
        [[
            {"RIGHT_ID": "verb", "RIGHT_ATTRS": {"POS": {"IN": ["VERB", "AUX"]}}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "subject", "RIGHT_ATTRS": {"DEP": {"IN": ["nsubj", "nsubjpass"]}}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "xcomp", "RIGHT_ATTRS": {"DEP": "xcomp"}},
        ]],
    )

    matcher.add(
        "PASSIVE_AGENT",
        [[
            {"RIGHT_ID": "verb", "RIGHT_ATTRS": {"POS": {"IN": ["VERB", "AUX"]}}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "patient", "RIGHT_ATTRS": {"DEP": "nsubjpass"}},
            {"LEFT_ID": "verb", "REL_OP": ">", "RIGHT_ID": "agent", "RIGHT_ATTRS": {"DEP": "agent"}},
            {"LEFT_ID": "agent", "REL_OP": ">", "RIGHT_ID": "actor", "RIGHT_ATTRS": {"DEP": "pobj"}},
        ]],
    )

    matcher.add(
        "APPOSITIVE",
        [[
            {"RIGHT_ID": "head", "RIGHT_ATTRS": {"POS": {"IN": ["NOUN", "PROPN", "PRON"]}}},
            {"LEFT_ID": "head", "REL_OP": ">", "RIGHT_ID": "appos", "RIGHT_ATTRS": {"DEP": "appos"}},
        ]],
    )

    matcher.add(
        "RELATIVE_CLAUSE",
        [[
            {"RIGHT_ID": "head", "RIGHT_ATTRS": {"POS": {"IN": ["NOUN", "PROPN", "PRON"]}}},
            {"LEFT_ID": "head", "REL_OP": ">", "RIGHT_ID": "relcl", "RIGHT_ATTRS": {"DEP": "relcl"}},
        ]],
    )

    return matcher


def project_sentence(sent: Span, coref_annotation: CorefAnnotation) -> List[Projection]:
    projections: List[Projection] = []
    noun_chunks = noun_chunk_map(sent)
    seen = set()

    for token in sent:
        for child in token.children:
            if child.dep_ == "appos":
                project_appositive(projections, seen, token, child, noun_chunks, coref_annotation)

    for token in sent:
        if token.pos_ in {"VERB", "AUX"}:
            subject = resolve_subject(token, noun_chunks, coref_annotation)

            for child in token.children:
                if child.dep_ in {"ccomp", "xcomp"} and subject:
                    if child.dep_ == "xcomp":
                        projection = Projection(
                            subject=subject,
                            relation=token.lemma_,
                            object=child.lemma_,
                            pattern="verb_xcomp_action",
                        )
                        add_projection(projections, seen, projection)
                    else:
                        child_subject = resolve_subject(child, noun_chunks, coref_annotation)
                        if child_subject:
                            projection = Projection(
                                subject=subject,
                                relation=token.lemma_,
                                object=child_subject,
                                pattern="verb_ccomp_subject",
                            )
                            add_projection(projections, seen, projection)

            for child in token.children:
                if child.dep_ in {"dobj", "obj"} and subject:
                    for object_token in with_conjuncts(child):
                        projection = Projection(
                            subject=subject,
                            relation=token.lemma_,
                            object=expand_entity_phrase(object_token, noun_chunks, coref_annotation),
                            pattern="verb_object",
                        )
                        add_projection(projections, seen, projection)

                if child.dep_ == "prep" and subject:
                    for prep_child in child.children:
                        if prep_child.dep_ == "pobj":
                            for object_token in with_conjuncts(prep_child):
                                projection = Projection(
                                    subject=subject,
                                    relation=f"{token.lemma_}_{child.lemma_}",
                                    object=expand_entity_phrase(object_token, noun_chunks, coref_annotation),
                                    pattern="verb_prep_object",
                                )
                                add_projection(projections, seen, projection)

                if child.dep_ == "advcl":
                    adv_subject = resolve_subject(child, noun_chunks, coref_annotation)
                    acomp = first_child(child, "acomp")
                    if adv_subject and acomp:
                        projection = Projection(
                            subject=adv_subject,
                            relation=child.lemma_,
                            object=expand_property_phrase(acomp),
                            pattern="advcl_property",
                        )
                        add_projection(projections, seen, projection)

            passive_patient = first_child(token, "nsubjpass")
            if passive_patient is not None:
                patient = expand_entity_phrase(passive_patient, noun_chunks, coref_annotation)
                for child in token.children:
                    if child.dep_ == "agent":
                        for agent_child in child.children:
                            if agent_child.dep_ == "pobj":
                                for actor_token in with_conjuncts(agent_child):
                                    projection = Projection(
                                        subject=expand_entity_phrase(actor_token, noun_chunks, coref_annotation),
                                        relation=token.lemma_,
                                        object=patient,
                                        pattern="passive_agent",
                                    )
                                    add_projection(projections, seen, projection)

            if token.lemma_ == "be":
                attr = first_child(token, "attr")
                if attr and subject:
                    projection = Projection(
                        subject=subject,
                        relation="is_a",
                        object=expand_entity_phrase(attr, noun_chunks, coref_annotation),
                        pattern="copula_attribute",
                    )
                    add_projection(projections, seen, projection)

                acomp = first_child(token, "acomp")
                if acomp and subject:
                    projection = Projection(
                        subject=subject,
                        relation="has_property",
                        object=expand_property_phrase(acomp),
                        pattern="copula_property",
                    )
                    add_projection(projections, seen, projection)

    return projections


def resolve_subject(token: Token, noun_chunks: Dict[int, Span], coref_annotation: CorefAnnotation) -> Optional[str]:
    direct = first_child(token, "nsubj") or first_child(token, "nsubjpass")
    if direct:
        if token.dep_ == "relcl" and direct.text.lower() in {"that", "which", "who", "whom"}:
            return expand_entity_phrase(token.head, noun_chunks, coref_annotation)
        return expand_entity_phrase(direct, noun_chunks, coref_annotation)

    if token.dep_ == "conj" and token.head != token:
        head_subject = resolve_subject(token.head, noun_chunks, coref_annotation)
        if head_subject:
            return head_subject

    if token.dep_ == "xcomp" and token.head != token:
        head_subject = resolve_subject(token.head, noun_chunks, coref_annotation)
        if head_subject:
            return head_subject

    return None


def first_child(token: Token, dep: str) -> Optional[Token]:
    for child in token.children:
        if child.dep_ == dep:
            return child
    return None


def noun_chunk_map(sent: Span) -> Dict[int, Span]:
    output: Dict[int, Span] = {}
    for chunk in sent.noun_chunks:
        output[chunk.root.i] = chunk
    return output


def expand_entity_phrase(
    token: Token,
    noun_chunks: Dict[int, Span],
    coref_annotation: CorefAnnotation,
    prefer_appos_identity: bool = True,
) -> str:
    if prefer_appos_identity:
        identity = appositive_identity(token, noun_chunks, coref_annotation)
        if identity is not None:
            return identity

    chunk = noun_chunks.get(token.i)
    if chunk is not None:
        anchor = resolve_span_anchor(chunk, coref_annotation)
        if anchor is not None:
            return anchor
        return chunk.text

    subtree_tokens = [tok for tok in token.subtree if tok.dep_ != "punct"]
    if not subtree_tokens:
        anchor = resolve_token_anchor(token, coref_annotation)
        return anchor if anchor is not None else token.text

    start = subtree_tokens[0].i
    end = subtree_tokens[-1].i + 1
    span = token.doc[start:end]
    anchor = resolve_span_anchor(span, coref_annotation)
    if anchor is not None:
        return anchor
    return span.text


def expand_property_phrase(token: Token) -> str:
    kept = [child for child in token.children if child.dep_ in {"advmod", "npadvmod"}]
    span_tokens = sorted([*kept, token], key=lambda tok: tok.i)
    start = span_tokens[0].i
    end = span_tokens[-1].i + 1
    return token.doc[start:end].text


def add_projection(projections: List[Projection], seen: set, projection: Projection) -> None:
    key = (projection.subject, projection.relation, projection.object, projection.pattern)
    if key in seen:
        return
    seen.add(key)
    projections.append(projection)


def resolve_span_anchor(span: Span, coref_annotation: CorefAnnotation) -> Optional[str]:
    if span.text.lower() in NON_REWRITE_PRONOUNS:
        return None
    char_span = (span.start_char, span.end_char)
    return coref_annotation.span_to_anchor.get(char_span)


def resolve_token_anchor(token: Token, coref_annotation: CorefAnnotation) -> Optional[str]:
    if token.text.lower() in NON_REWRITE_PRONOUNS:
        return None
    char_span = (token.idx, token.idx + len(token.text))
    return coref_annotation.span_to_anchor.get(char_span)

def with_conjuncts(token: Token) -> List[Token]:
    return [token, *list(token.conjuncts)]


def appositive_identity(
    token: Token,
    noun_chunks: Dict[int, Span],
    coref_annotation: CorefAnnotation,
) -> Optional[str]:
    if token.pos_ == "PROPN":
        return None

    for child in token.children:
        if child.dep_ == "appos" and contains_proper_noun(child):
            return expand_entity_phrase(child, noun_chunks, coref_annotation, prefer_appos_identity=False)
    return None


def contains_proper_noun(token: Token) -> bool:
    return token.pos_ == "PROPN" or any(child.pos_ == "PROPN" for child in token.subtree)


def project_appositive(
    projections: List[Projection],
    seen: set,
    head: Token,
    appos: Token,
    noun_chunks: Dict[int, Span],
    coref_annotation: CorefAnnotation,
) -> None:
    head_surface = expand_entity_phrase(head, noun_chunks, coref_annotation, prefer_appos_identity=False)
    appos_surface = expand_entity_phrase(appos, noun_chunks, coref_annotation, prefer_appos_identity=False)

    if contains_proper_noun(appos) and head.pos_ != "PROPN":
        projection = Projection(
            subject=head_surface,
            relation="identity",
            object=appos_surface,
            pattern="appos_identity",
        )
        add_projection(projections, seen, projection)
        return

    if head.pos_ == "PROPN" and appos.pos_ in {"NOUN", "PROPN"}:
        projection = Projection(
            subject=head_surface,
            relation="is_a",
            object=appos_surface,
            pattern="appos_is_a",
        )
        add_projection(projections, seen, projection)


if __name__ == "__main__":
    raise SystemExit(main())
