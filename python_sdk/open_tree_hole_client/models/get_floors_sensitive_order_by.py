from enum import Enum


class GetFloorsSensitiveOrderBy(str, Enum):
    TIME_CREATED = "time_created"
    TIME_UPDATED = "time_updated"

    def __str__(self) -> str:
        return str(self.value)
