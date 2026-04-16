class FaturaError(Exception):
    """Base exception for the project."""


class ConfigError(FaturaError):
    """Invalid or missing configuration."""


class PortalError(FaturaError):
    """Base for all portal interaction errors."""

    def __init__(self, message: str, uc: str | None = None):
        self.uc = uc
        super().__init__(message)


class LoginError(PortalError):
    """Authentication failed."""


class CaptchaError(PortalError):
    """CAPTCHA detected and could not be bypassed."""


class SessionExpiredError(PortalError):
    """Session timed out mid-operation."""


class DownloadError(PortalError):
    """PDF download failed or timed out."""


class LayoutChangedError(PortalError):
    """Expected DOM elements not found -- portal layout probably changed."""


class RateLimitError(PortalError):
    """Too many requests to portal."""


class ParserError(FaturaError):
    """PDF could not be parsed."""


class FieldNotFoundError(ParserError):
    """A required field was not found in the PDF text."""

    def __init__(self, field: str, context: str = ""):
        self.field = field
        self.context = context
        super().__init__(f"Campo '{field}' não encontrado no PDF. Contexto: {context[:200]}")


class RepositoryError(FaturaError):
    """Database operation failed."""
