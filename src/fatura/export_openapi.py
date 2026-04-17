from __future__ import annotations

import json
from pathlib import Path

from fatura.api import create_app
from fatura.config import AppConfig


def main() -> None:
    app = create_app(config=AppConfig())
    output = Path("docs/openapi.json")
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(json.dumps(app.openapi(), ensure_ascii=False, indent=2), encoding="utf-8")
    print(output)


if __name__ == "__main__":
    main()
