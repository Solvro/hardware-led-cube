import numpy as np
import numpy.typing as npt
from abc import ABC, abstractmethod
import json


def rgbToColor(r: int, g: int, b: int) -> int:
    return (r << 16) | (g << 8) | b


class StateParser(ABC):
    @abstractmethod
    def parse(self, leds: npt.NDArray) -> str:
        ...


class PrototypeJsonifier(StateParser):
    def __init__(self) -> None:
        self.previous_states: list[npt.NDArray] = []

    def parse(self, leds: npt.NDArray) -> str:
        data: dict = {
            "id": f"Frame_{len(self.previous_states)}",
            "state": leds.tolist()
        }

        json_string: str = json.dumps(data, indent=4)
        return json_string
