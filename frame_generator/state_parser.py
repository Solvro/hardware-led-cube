import numpy.typing as npt
from abc import ABC, abstractmethod
import json


def rgbToColor(r: int, g: int, b: int) -> int:
    return (int(r) << 16) | (int(g) << 8) | int(b)


class StateParser(ABC):
    @abstractmethod
    def handle_frame(self, leds: npt.NDArray) -> None:
        ...

    @property
    @abstractmethod
    def extension(self) -> str:
        ...

    @abstractmethod
    def get_parsed_results(self) -> str:
        ...


class PrototypeJsonifier(StateParser):
    def __init__(self) -> None:
        self.previous_states: list[list[list[list[int]]]] = []

    def handle_frame(self, leds: npt.NDArray) -> None:
        color_state = leds.reshape(
            (leds.shape[0], leds.shape[1], leds.shape[2], 3))

        int_color_state: list[list[list[int]]] = [[[rgbToColor(r, g, b) for r, g, b in row]
                                                   for row in plane] for plane in color_state]

        self.previous_states.append(int_color_state)

    def get_parsed_results(self) -> str:
        json_string: str = json.dumps(self.previous_states, indent=4)
        return json_string

    @property
    def extension(self) -> str:
        return "json"
