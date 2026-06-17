"""Self-contained branching functions + a class, for tsma Python dogfood."""


def classify(n):
    if n > 0:
        return "pos"
    return "nonpos"


def absval(n):
    if n < 0:
        return -n
    return n


async def fetch(value):
    return value


class Calculator:
    def add(self, a, b):
        return a + b

    def sign(self, n):
        if n > 0:
            return 1
        elif n < 0:
            return -1
        return 0

    @staticmethod
    def mul(a, b):
        def _double(x):
            return x + x

        return _double(a * b) - a * b


def _private():
    return 0
