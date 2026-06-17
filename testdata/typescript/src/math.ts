// Self-contained, branch-bearing functions for D1/D2/D3 fixtures.

export function add(a: number, b: number): number {
  if (a > b) {
    return a + b;
  }
  return b + a;
}

export const classify = (n: number): string => {
  if (n < 0) {
    return "negative";
  }
  if (n === 0) {
    return "zero";
  }
  return "positive";
};

function internalHelper(x: number): number {
  return x * 2;
}

export function double(x: number): number {
  return internalHelper(x);
}
