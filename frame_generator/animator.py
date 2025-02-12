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

    def stop_function(self, leds: npt.NDArray) -> bool:
        return False


class SampleAnimator(Animator):
    def internal_setup(self) -> None:
        self.current_x_index = 0

    def start_function(self, leds) -> npt.NDArray:
        return self.update_function(leds)

    def update_function(self, leds) -> npt.NDArray:
        for x in range(leds.shape[0]):
            for y in range(leds.shape[1]):
                for z in range(leds.shape[2]):
                    if x == self.current_x_index:
                        leds[x][y][z] = [1, 1, 1]
                    else:
                        leds[x][y][z] = [0, 0, 0]

        self.current_x_index += 1

        return leds

    def stop_function(self, leds) -> bool:
        return self.current_x_index == leds.shape[0]
