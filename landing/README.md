# Landing Website

Lean static landing site intended for internal hosting.

## Run locally

```bash
cd landing
python3 -m http.server 8080 --bind 0.0.0.0
```

Then open:

- `http://localhost:8080` on the Mac mini
- `http://<mac-mini-lan-ip>:8080` from other devices on your network

## Edit founder profiles

Founder cards are data-driven from `founders.js`.

Add another founder by appending an object to `window.founders`:

```js
{
  name: "Name",
  role: "Founder",
  bio: "Short biography..."
}
```

## Update the OctOS image

The hero image is at `assets/octos-mark.svg`.

To use your own photo/logo, replace that file or update the `<img src>` in
`index.html`.
