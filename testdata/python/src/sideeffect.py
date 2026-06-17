"""A module that runs code at import time (import side-effect case).

ast.parse never executes this, so D1 indexing is unaffected; the D5 loop must
import it for measurement, which is what makes this a useful fixture.
"""

_LOADED = []


def _register():
    _LOADED.append("loaded")


# Side-effect: this runs the moment the module is imported.
_register()


def greet(name):
    if name:
        return "hi " + name
    return "hi"
