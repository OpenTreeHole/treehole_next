from enum import Enum


class ModelsMessageType(str, Enum):
    FAVORITE = "favorite"
    MAIL = "mail"
    MENTION = "mention"
    MODIFY = "modify"
    PERMISSION = "permission"
    REPLY = "reply"
    REPORT = "report"
    REPORT_DEALT = "report_dealt"
    SENSITIVE = "sensitive"

    def __str__(self) -> str:
        return str(self.value)
