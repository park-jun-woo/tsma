// One file, multiple functions: a class with methods + a top-level function +
// a const arrow, with multi-line signatures the line-based indexer mangles.

export class Rectangle {
  private width: number;
  private height: number;

  constructor(width: number, height: number) {
    this.width = width;
    this.height = height;
  }

  Area(): number {
    return this.width * this.height;
  }

  scale(
    factor: number,
  ): Rectangle {
    return new Rectangle(this.width * factor, this.height * factor);
  }
}

export function makeSquare(side: number): Rectangle {
  return new Rectangle(side, side);
}

export const perimeter = (r: Rectangle): number => {
  return 0;
};
