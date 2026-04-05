import edu.stanford.nlp.ie.util.RelationTriple;
import edu.stanford.nlp.ling.CoreAnnotations;
import edu.stanford.nlp.ling.CoreLabel;
import edu.stanford.nlp.naturalli.NaturalLogicAnnotations;
import edu.stanford.nlp.pipeline.Annotation;
import edu.stanford.nlp.pipeline.StanfordCoreNLP;
import edu.stanford.nlp.util.CoreMap;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashSet;
import java.util.List;
import java.util.Properties;
import java.util.Set;

public class StructuralProjectionDemo {
  public static void main(String[] args) {
    Config config = Config.parse(args);
    ThemeProfile theme = ThemeProfile.load(config.themeId, config.themeNote);
    String normalizedInput = preprocessInput(config.text);

    Properties props = new Properties();
    props.setProperty("annotators", "tokenize,ssplit,pos,lemma,ner,depparse,natlog,openie");

    StanfordCoreNLP pipeline = new StanfordCoreNLP(props);
    Annotation doc = new Annotation(normalizedInput);
    pipeline.annotate(doc);

    System.out.println("INPUT");
    System.out.println(config.text);
    System.out.println();

    System.out.println("MODE");
    System.out.println("  " + (config.rawMode ? "raw" : "normalized"));
    System.out.println();

    if (!config.text.equals(normalizedInput)) {
      System.out.println("PREPROCESSED");
      System.out.println(normalizedInput);
      System.out.println();
    }

    System.out.println("THEME");
    System.out.println("  id: " + theme.id);
    System.out.println("  description: " + theme.description);
    if (!config.themeNote.isBlank()) {
      System.out.println("  theme_note: " + config.themeNote);
    }
    System.out.println();

    List<CoreMap> sentences = doc.get(CoreAnnotations.SentencesAnnotation.class);
    List<EntityResolution> chunkResolvedEntities = new ArrayList<>();
    for (int i = 0; i < sentences.size(); i++) {
      CoreMap sentence = sentences.get(i);
      String sentenceText = sentence.toString();

      System.out.println("SENTENCE " + (i + 1));
      System.out.println(sentenceText);

      List<String> mentionSurfaces = extractMentionSurfaces(sentence);
      List<EntityResolution> sentenceMentionResolutions = new ArrayList<>();
      if (!mentionSurfaces.isEmpty()) {
        System.out.println("  MENTION CANDIDATES");
        for (String mention : mentionSurfaces) {
          EntityResolution mentionResolution = normalizeEntity(
            mention,
            sentenceText,
            sentenceText,
            chunkResolvedEntities,
            theme,
            config.rawMode
          );
          sentenceMentionResolutions.add(mentionResolution);
          System.out.println("    " + mention + " -> " + mentionResolution.id + " [" + mentionResolution.label + "] score=" + fmt(mentionResolution.score));
        }
      }

      Collection<RelationTriple> triples = sentence.get(NaturalLogicAnnotations.RelationTriplesAnnotation.class);
      if (triples == null || triples.isEmpty()) {
        System.out.println("  No triples");
        System.out.println();
        continue;
      }

      int tripleIndex = 1;
      for (RelationTriple triple : triples) {
        RawTriple raw = new RawTriple(
          triple.subjectGloss(),
          triple.relationGloss(),
          triple.objectGloss(),
          triple.confidence
        );

        String tripleContext = raw.subject + " " + raw.relation + " " + raw.object;
        EntityResolution subject = normalizeEntity(raw.subject, tripleContext, sentenceText, chunkResolvedEntities, theme, config.rawMode);
        chunkResolvedEntities.add(subject);

        EntityResolution object = normalizeEntity(raw.object, tripleContext, sentenceText, chunkResolvedEntities, theme, config.rawMode);
        chunkResolvedEntities.add(object);

        RelationshipResolution relation = normalizeRelationship(raw.relation);

        System.out.println("  RAW TRIPLE " + tripleIndex);
        System.out.println("    subject: " + raw.subject);
        System.out.println("    relation: " + raw.relation);
        System.out.println("    object: " + raw.object);
        System.out.println("    confidence: " + raw.confidence);

        System.out.println("  NORMALIZED TRIPLE " + tripleIndex);
        System.out.println("    subject: " + subject.id + " [" + subject.label + "] score=" + fmt(subject.score));
        System.out.println("    relation: " + relation.id + " [" + relation.label + "] score=" + fmt(relation.score));
        System.out.println("    object: " + object.id + " [" + object.label + "] score=" + fmt(object.score));
        tripleIndex++;
      }

      chunkResolvedEntities.addAll(sentenceMentionResolutions);

      System.out.println();
    }
  }

  private static List<String> extractMentionSurfaces(CoreMap sentence) {
    List<CoreLabel> tokens = sentence.get(CoreAnnotations.TokensAnnotation.class);
    List<String> mentions = new ArrayList<>();

    StringBuilder currentNer = new StringBuilder();
    String currentNerTag = "O";

    for (CoreLabel token : tokens) {
      String text = token.word();
      String ner = token.ner();

      if (ner != null && !"O".equals(ner)) {
        if (ner.equals(currentNerTag) && currentNer.length() > 0) {
          currentNer.append(" ").append(text);
        } else {
          flushMention(currentNer, mentions);
          currentNer = new StringBuilder(text);
          currentNerTag = ner;
        }
      } else {
        flushMention(currentNer, mentions);
        currentNerTag = "O";
      }
    }
    flushMention(currentNer, mentions);

    StringBuilder nounPhrase = new StringBuilder();
    boolean hasNoun = false;
    for (CoreLabel token : tokens) {
      String word = token.word();
      String pos = token.tag();

      boolean keep = pos != null && (
        pos.startsWith("NN")
          || pos.startsWith("JJ")
      );

      if (keep) {
        if (nounPhrase.length() > 0) {
          nounPhrase.append(" ");
        }
        nounPhrase.append(word);
        if (pos.startsWith("NN")) {
          hasNoun = true;
        }
      } else {
        if (nounPhrase.length() > 0 && hasNoun) {
          mentions.add(nounPhrase.toString());
        }
        nounPhrase = new StringBuilder();
        hasNoun = false;
      }
    }
    if (nounPhrase.length() > 0 && hasNoun) {
      mentions.add(nounPhrase.toString());
    }

    List<String> deduped = new ArrayList<>();
    Set<String> seen = new HashSet<>();
    for (String mention : mentions) {
      String normalized = normalizeText(mention);
      if (normalized.isBlank() || isPronoun(normalized) || isSelfPronoun(normalized)) {
        continue;
      }
      if (seen.add(normalized)) {
        deduped.add(mention);
      }
    }
    return deduped;
  }

  private static void flushMention(StringBuilder builder, List<String> mentions) {
    if (builder.length() > 0) {
      mentions.add(builder.toString());
      builder.setLength(0);
    }
  }

  private static EntityResolution normalizeEntity(
    String surface,
    String tripleContext,
    String sentenceText,
    List<EntityResolution> localResolvedEntities,
    ThemeProfile theme,
    boolean rawMode
  ) {
    if (rawMode) {
      return EntityResolution.fromSurface(surface);
    }

    String normalizedSurface = normalizeText(surface);
    if (normalizedSurface.isBlank()) {
      return EntityResolution.provisional(surface);
    }

    if (isAffectWord(normalizedSurface)) {
      return EntityResolution.provisional(surface);
    }

    if (isSelfPronoun(normalizedSurface)) {
      CanonicalEntity self = theme.findEntity("self");
      return EntityResolution.fromCanonical(surface, self, 1.0);
    }

    if (isPronoun(normalizedSurface)) {
      EntityResolution localWinner = pickLocalPronounTarget(surface, tripleContext, sentenceText, localResolvedEntities, theme);
      if (localWinner != null) {
        return localWinner.reboundFrom(surface);
      }
    }

    CanonicalEntity best = null;
    double bestScore = Double.NEGATIVE_INFINITY;
    for (CanonicalEntity entity : theme.entities) {
      double score = scoreEntity(entity, normalizedSurface, tripleContext, sentenceText, theme);
      if (score > bestScore) {
        bestScore = score;
        best = entity;
      }
    }

    if (best == null || bestScore < 2.0) {
      return EntityResolution.provisional(surface);
    }

    return EntityResolution.fromCanonical(surface, best, bestScore);
  }

  private static EntityResolution pickLocalPronounTarget(
    String surface,
    String tripleContext,
    String sentenceText,
    List<EntityResolution> localResolvedEntities,
    ThemeProfile theme
  ) {
    EntityResolution best = null;
    double bestScore = Double.NEGATIVE_INFINITY;
    int count = localResolvedEntities.size();

    for (int i = 0; i < count; i++) {
      EntityResolution candidate = localResolvedEntities.get(i);
      if (candidate.isProvisional) {
        continue;
      }
      double recencyBoost = 1.5 * ((double) (i + 1) / count);
      double lexicalContext = tokenOverlapScore(
        tokenSet(candidate.label + " " + String.join(" ", candidate.keywords)),
        tokenSet(tripleContext + " " + sentenceText)
      );
      double activeBoost = theme.activeEntityIds.contains(candidate.id) ? 0.75 : 0.0;
      double score = recencyBoost + lexicalContext + activeBoost;
      if (score > bestScore) {
        bestScore = score;
        best = candidate;
      }
    }

    if (best == null) {
      return null;
    }

    return best.reboundFrom(surface, bestScore);
  }

  private static double scoreEntity(
    CanonicalEntity entity,
    String normalizedSurface,
    String tripleContext,
    String sentenceText,
    ThemeProfile theme
  ) {
    double aliasScore = bestLexicalScore(normalizedSurface, entity.allNames());
    double contextScore = tokenOverlapScore(
      tokenSet(tripleContext + " " + sentenceText),
      tokenSet(entity.label + " " + String.join(" ", entity.keywords))
    );
    double themeScore = tokenOverlapScore(theme.themeTokens, tokenSet(entity.label + " " + String.join(" ", entity.keywords)));
    double activeBoost = theme.activeEntityIds.contains(entity.id) ? 0.75 : 0.0;
    return aliasScore + contextScore + themeScore + activeBoost;
  }

  private static RelationshipResolution normalizeRelationship(String surface) {
    String normalized = normalizeText(surface);
    if (normalized.isBlank()) {
      return RelationshipResolution.provisional(surface);
    }
    return RelationshipResolution.fromSurface(surface, normalized);
  }

  private static boolean isSelfPronoun(String text) {
    return Set.of("i", "me", "my", "myself").contains(text);
  }

  private static boolean isPronoun(String text) {
    return Set.of("it", "this", "that", "they", "them", "he", "she", "him", "her", "you", "u").contains(text);
  }

  private static boolean isAffectWord(String text) {
    String norm = normalizeText(text);
    return Set.of("frustrating", "draining", "good", "bad", "stressful", "exciting").contains(norm);
  }

  private static double bestLexicalScore(String surface, List<String> candidates) {
    double best = 0.0;
    for (String candidate : candidates) {
      String normalizedCandidate = normalizeText(candidate);
      if (normalizedCandidate.equals(surface)) {
        best = Math.max(best, 3.0);
      } else if (normalizedCandidate.contains(surface) || surface.contains(normalizedCandidate)) {
        best = Math.max(best, 2.0);
      } else {
        best = Math.max(best, tokenOverlapScore(tokenSet(surface), tokenSet(normalizedCandidate)));
      }
    }
    return best;
  }

  private static double tokenOverlapScore(Set<String> left, Set<String> right) {
    if (left.isEmpty() || right.isEmpty()) {
      return 0.0;
    }
    Set<String> intersection = new HashSet<>(left);
    intersection.retainAll(right);
    return (double) intersection.size() / Math.max(left.size(), right.size());
  }

  private static Set<String> tokenSet(String text) {
    Set<String> output = new HashSet<>();
    for (String part : normalizeText(text).split(" ")) {
      if (!part.isBlank()) {
        output.add(part);
      }
    }
    return output;
  }

  private static String normalizeText(String text) {
    return text.toLowerCase()
      .replaceAll("[^a-z0-9 ]", " ")
      .replaceAll("\\s+", " ")
      .trim();
  }

  private static String preprocessInput(String text) {
    String output = text;

    output = output.replaceAll("(?i)\\bu\\b", "you");
    output = output.replaceAll("(?i)\\bits\\b", "it is");
    output = output.replaceAll("(?i)\\bim\\b", "I am");
    output = output.replaceAll("(?i)\\bdont\\b", "do not");
    output = output.replaceAll("(?i)\\bcant\\b", "can not");
    output = output.replaceAll("(?i)\\bwont\\b", "will not");

    return output;
  }

  private static String fmt(double value) {
    return String.format("%.2f", value);
  }

  private static final class Config {
    final String themeId;
    final String themeNote;
    final boolean rawMode;
    final String text;

    private Config(String themeId, String themeNote, boolean rawMode, String text) {
      this.themeId = themeId;
      this.themeNote = themeNote;
      this.rawMode = rawMode;
      this.text = text;
    }

    static Config parse(String[] args) {
      String themeId = "software_work";
      String themeNote = "";
      boolean rawMode = false;
      List<String> textParts = new ArrayList<>();

      for (int i = 0; i < args.length; i++) {
        String arg = args[i];
        if ("--theme".equals(arg) && i + 1 < args.length) {
          themeId = args[++i];
        } else if ("--theme-note".equals(arg) && i + 1 < args.length) {
          themeNote = args[++i];
        } else if ("--raw".equals(arg)) {
          rawMode = true;
        } else {
          textParts.add(arg);
        }
      }

      String text = textParts.isEmpty()
        ? "I am outside doing construction again today."
        : String.join(" ", textParts);

      return new Config(themeId, themeNote, rawMode, text);
    }
  }

  private static final class RawTriple {
    final String subject;
    final String relation;
    final String object;
    final double confidence;

    private RawTriple(String subject, String relation, String object, double confidence) {
      this.subject = subject;
      this.relation = relation;
      this.object = object;
      this.confidence = confidence;
    }
  }

  private static final class EntityResolution {
    final String surfaceForm;
    final String id;
    final String label;
    final List<String> keywords;
    final double score;
    final boolean isProvisional;

    private EntityResolution(String surfaceForm, String id, String label, List<String> keywords, double score, boolean isProvisional) {
      this.surfaceForm = surfaceForm;
      this.id = id;
      this.label = label;
      this.keywords = keywords;
      this.score = score;
      this.isProvisional = isProvisional;
    }

    static EntityResolution fromCanonical(String surfaceForm, CanonicalEntity entity, double score) {
      return new EntityResolution(surfaceForm, entity.id, entity.label, entity.keywords, score, false);
    }

    static EntityResolution fromSurface(String surfaceForm) {
      String normalized = normalizeText(surfaceForm);
      String id = normalized.isBlank() ? "entity:unknown" : "entity:" + normalized.replace(' ', '_');
      return new EntityResolution(surfaceForm, id, surfaceForm, List.of(), 1.0, false);
    }

    static EntityResolution provisional(String surfaceForm) {
      String normalized = normalizeText(surfaceForm);
      String id = normalized.isBlank() ? "entity:unknown" : "entity:" + normalized.replace(' ', '_');
      return new EntityResolution(surfaceForm, id, surfaceForm, List.of(), 0.0, true);
    }

    EntityResolution reboundFrom(String newSurface) {
      return new EntityResolution(newSurface, id, label, keywords, score, isProvisional);
    }

    EntityResolution reboundFrom(String newSurface, double newScore) {
      return new EntityResolution(newSurface, id, label, keywords, newScore, isProvisional);
    }
  }

  private static final class RelationshipResolution {
    final String surfaceForm;
    final String id;
    final String label;
    final double score;
    final boolean isProvisional;

    private RelationshipResolution(String surfaceForm, String id, String label, double score, boolean isProvisional) {
      this.surfaceForm = surfaceForm;
      this.id = id;
      this.label = label;
      this.score = score;
      this.isProvisional = isProvisional;
    }

    static RelationshipResolution fromSurface(String surfaceForm, String normalizedSurface) {
      String id = "rel:" + normalizedSurface.replace(' ', '_');
      return new RelationshipResolution(surfaceForm, id, normalizedSurface, 1.0, false);
    }

    static RelationshipResolution provisional(String surfaceForm) {
      String normalized = normalizeText(surfaceForm);
      String id = normalized.isBlank() ? "rel:unknown" : "rel:" + normalized.replace(' ', '_');
      return new RelationshipResolution(surfaceForm, id, surfaceForm, 0.0, true);
    }
  }

  private static final class CanonicalEntity {
    final String id;
    final String label;
    final List<String> aliases;
    final List<String> keywords;

    private CanonicalEntity(String id, String label, List<String> aliases, List<String> keywords) {
      this.id = id;
      this.label = label;
      this.aliases = aliases;
      this.keywords = keywords;
    }

    List<String> allNames() {
      List<String> names = new ArrayList<>();
      names.add(label);
      names.addAll(aliases);
      return names;
    }
  }

  private static final class ThemeProfile {
    final String id;
    final String description;
    final List<CanonicalEntity> entities;
    final Set<String> activeEntityIds;
    final Set<String> themeTokens;

    private ThemeProfile(
      String id,
      String description,
      List<CanonicalEntity> entities,
      Set<String> activeEntityIds,
      Set<String> themeTokens
    ) {
      this.id = id;
      this.description = description;
      this.entities = entities;
      this.activeEntityIds = activeEntityIds;
      this.themeTokens = themeTokens;
    }

    CanonicalEntity findEntity(String id) {
      return entities.stream()
        .filter(entity -> entity.id.equals(id))
        .findFirst()
        .orElseThrow();
    }

    static ThemeProfile load(String id, String themeNote) {
      if (!"software_work".equals(id)) {
        throw new IllegalArgumentException("Unknown theme: " + id + ". Only software_work is implemented right now.");
      }

      List<CanonicalEntity> entities = List.of(
        new CanonicalEntity("self", "self", List.of("i", "me", "myself"), List.of("person", "user")),
        new CanonicalEntity("assistant", "assistant", List.of("you", "u"), List.of("helper", "system", "assistant")),
        new CanonicalEntity("architecture", "architecture", List.of("system architecture", "design"), List.of("software", "system", "design", "structure")),
        new CanonicalEntity("library", "library", List.of("new library", "dependency", "package"), List.of("tooling", "dependency", "package", "framework")),
        new CanonicalEntity("terraform", "terraform", List.of("Terraform"), List.of("iac", "infrastructure", "tool", "language", "software")),
        new CanonicalEntity("language", "language", List.of("new language", "programming language"), List.of("syntax", "tool", "software", "language")),
        new CanonicalEntity("project", "project", List.of("this project", "app", "system"), List.of("work", "software", "build")),
        new CanonicalEntity("toolset", "toolset", List.of("tools", "tooling"), List.of("tools", "stack", "workflow")),
        new CanonicalEntity("work", "work", List.of("job"), List.of("task", "career", "effort")),
        new CanonicalEntity("outside", "outside", List.of("outdoors"), List.of("environment", "location")),
        new CanonicalEntity("construction", "construction", List.of("building work"), List.of("physical", "work", "labor")),
        new CanonicalEntity("device", "device", List.of("hardware"), List.of("machine", "computer")),
        new CanonicalEntity("tv", "tv", List.of("television", "roku tv"), List.of("screen", "display", "device")),
        new CanonicalEntity("pi", "pi", List.of("raspberry pi", "pi zero"), List.of("device", "probe", "controller"))
      );

      Set<String> activeEntityIds = Set.of("self", "assistant", "architecture", "library", "terraform", "language", "project", "toolset", "work");
      Set<String> themeTokens = tokenSet("software work architecture libraries tooling coding systems terraform infrastructure language " + themeNote);

      return new ThemeProfile(
        id,
        "Software/work normalization with active architecture and tooling bias.",
        entities,
        activeEntityIds,
        themeTokens
      );
    }
  }
}
