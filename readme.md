SVG/PNG of all pages can be obtained from `https://musescore.com/api/jmuse?id=6102579&index=[PAGE_NO, 0 indexed]&type=img&v2` but it requires authorization header. Authorization can be obtained from a js chunk in the file html.

Steps:

1. GET the main page (url provided as input)
2. Find the js chunk `<link href="https://musescore.com/static/public/build/musescore/202401/2946.ad0ca593092255f4a3f5b0e492146221.js" rel="preload" as="script">` yet to figure out how to filter it out
3. Use a regex match to get the Authorization header (TODO: find a better way)
4. Hit `https://musescore.com/api/jmuse?id=6102579&index=[PAGE_NO, 0-indexed]&type=img&v2` to get all pages
5. Generate html from svg/png files
6. Generate pdf from html (Still stuck on this one)

Usage:

```
./musescoreDownloader [url]
```
OR
```
go run . dl [url]
```
