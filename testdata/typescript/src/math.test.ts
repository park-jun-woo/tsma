// D2 content-based matching fixture: this test calls add() and classify() by
// name, so the content matcher must attribute it to those functions — not just
// to everything in math.ts by filename.
import { add, classify } from "./math";

describe("add", () => {
  it("adds the larger first", () => {
    expect(add(2, 1)).toBe(3);
    expect(add(1, 2)).toBe(3);
  });
});

describe("classify", () => {
  it("classifies sign", () => {
    expect(classify(-1)).toBe("negative");
    expect(classify(0)).toBe("zero");
    expect(classify(5)).toBe("positive");
  });
});
