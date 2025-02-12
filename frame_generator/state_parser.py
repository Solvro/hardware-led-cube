import numpy.typing as npt
from abc import ABC, abstractmethod
import json


def rgbToColor(r: int, g: int, b: int) -> int:
    return (int(r) << 16) | (int(g) << 8) | int(b)


class StateParser(ABC):
    @abstractmethod
    def parse(self, leds: npt.NDArray) -> str:
        ...


class PrototypeJsonifier(StateParser):
    def __init__(self) -> None:
        self.previous_states: list[npt.NDArray] = []

    def parse(self, leds: npt.NDArray) -> str:
        color_state = leds.reshape(
            (leds.shape[0], leds.shape[1], leds.shape[2], 3)).tolist()

        data: dict = {
            "id": f"Frame_{len(self.previous_states)}",
            "state": color_state
        }

        json_string: str = json.dumps(data, indent=4)
        return json_string
