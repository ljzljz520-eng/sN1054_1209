import { mkdir, writeFile } from "node:fs/promises";

await mkdir("dist", { recursive: true });
const page = `<!doctype html><html><head><meta charset="utf-8"><title>Family Itinerary Advisor</title></head><body><main><h1>Family Itinerary Advisor</h1><p>Plan a calm, child-friendly trip with a travel advisor.</p></main></body></html>`;
await writeFile("dist/index.html", page);
