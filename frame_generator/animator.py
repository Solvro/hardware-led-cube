from abc import ABC, abstractmethod
import numpy.typing as npt
from PIL import Image


class Animator(ABC):
    def __init__(self, cube_width: int) -> None:
        self.internal_setup(cube_width)

    @abstractmethod
    def internal_setup(self, cube_width: int) -> None:
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
    def internal_setup(self, cube_width: int) -> None:
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


class StaticSolvro(Animator):
    def internal_setup(self, cube_width: int) -> None:
        img = Image.open('frame_generator/solvro.png')
        img = img.convert("RGBA")

        img = img.resize((cube_width, cube_width))
        data = img.getdata()

        # Replace transparent pixels (alpha == 0) with black
        processed_data = [
            (0, 0, 0) if color[3] == 0 else color[:3]
            for color in data
        ]

        img.show()

        self.data = processed_data
        self.counter = 0

    def start_function(self, leds):
        return self.update_function(leds)

    def update_function(self, leds):
        self.counter += 1

        for x in range(leds.shape[0]):
            for y in range(leds.shape[1]):
                for z in range(leds.shape[2]):
                    leds[x][y][z] = self.data[x + y]

        if self.counter % 2 == 0:
            return leds

        return leds

    def stop_function(self, leds, frames):
        return frames > 0
