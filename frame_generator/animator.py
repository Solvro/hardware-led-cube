from abc import ABC, abstractmethod
from PIL.ImageFile import ImageFile
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


class TwoFacedSolvro(Animator):
    DISPLAY_PLANE = 4

    def internal_setup(self, cube_width: int) -> None:
        img: ImageFile = Image.open('frame_generator/solvro8x8x8.png')

        self.data = TwoFacedSolvro.get_formatted_image_data(img)
        self.counter = 0

    @staticmethod
    def get_formatted_image_data(img) -> list[list[float]]:
        data = list(img.getdata())
        data.reverse()
        data = [[c/255 for c in color[:3]] for color in data]
        return data

    def start_function(self, leds):
        return self.update_function(leds)

    def update_function(self, leds):
        self.counter += 1
        is_even_frame = self.counter % 2 == 0

        for x in range(leds.shape[0]):
            for y in range(leds.shape[1]):
                for z in range(leds.shape[2]):
                    leds[x][y][z] = [0, 0, 0]

                    if is_even_frame and z == self.DISPLAY_PLANE:
                        leds[x][y][z] = self.data[x + y * leds.shape[0]]
                    elif not is_even_frame and x == self.DISPLAY_PLANE:
                        leds[x][y][z] = self.data[z + y * leds.shape[0]]

        return leds

    def stop_function(self, leds, frames):
        return self.counter > 0
