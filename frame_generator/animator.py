from abc import ABC, abstractmethod
import numpy.typing as npt


class Animator(ABC):
    def __init__(self) -> None:
        self.internal_setup()

    @abstractmethod
    def internal_setup(self) -> None:
        ...

    @abstractmethod
    def start_function(self, leds: npt.NDArray) -> npt.NDArray:
        ...

    @abstractmethod
    def update_function(self, leds: npt.NDArray) -> npt.NDArray:
        ...

    def stop_function(self, leds: npt.NDArray, frames) -> bool:
        return False


class RGBMovingRight(Animator):
    def internal_setup(self) -> None:
        self.current_x_index = 0

    def start_function(self, leds) -> npt.NDArray:
        return self.update_function(leds)

    def update_function(self, leds) -> npt.NDArray:
        for x in range(leds.shape[0]):
            color = [1 if x == self.current_x_index else 0, 1 if x - 1 ==
                     self.current_x_index else 0, 1 if x - 2 == self.current_x_index else 0]
            for y in range(leds.shape[1]):
                for z in range(leds.shape[2]):
                    leds[x][y][z] = color

        self.current_x_index += 1

        return leds

    def stop_function(self, leds, frames) -> bool:
        return frames == leds.shape[0] - 1
