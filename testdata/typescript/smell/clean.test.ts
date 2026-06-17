// D4 negative fixture: legitimate idioms that must NOT fire (false-positive zero).
import { Rectangle, makeSquare } from "../src/shapes";

describe("clean", () => {
  it("uses public API and normal casts", () => {
    const r = makeSquare(4) as Rectangle; // legit `as T` cast — not `as any`
    const area = r.Area();
    const keys = Object.keys(r); // public iteration — not getOwnPropertyNames
    const vals = Object.values(r);
    // The word Reflect and "as any" appear in this comment but are not nodes.
    expect(area).toBe(16);
    expect(keys.length + vals.length).toBeGreaterThanOrEqual(0);
  });
});
