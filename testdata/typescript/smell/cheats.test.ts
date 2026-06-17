// D4 positive fixture: every escape hatch the TS smell rules must catch.
import { Rectangle } from "../src/shapes";

describe("cheats", () => {
  it("reaches into privates", () => {
    const r = new Rectangle(2, 3);
    const w = (r as any).width; // TS-REFL-TS-001
    const h = (r as any).getHeight(); // TS-REFL-TS-001
    const got = Reflect.get(r, "height"); // TS-REFL-TS-002
    Reflect.set(r, "width", 9); // TS-REFL-TS-002
    const names = Object.getOwnPropertyNames(r); // TS-REFL-TS-003
    const d = Object.getOwnPropertyDescriptor(r, "width"); // TS-REFL-TS-003
    expect(w + h + got + names.length + (d ? 1 : 0)).toBeGreaterThan(0);
  });
});
