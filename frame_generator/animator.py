from abc import ABC, abstractmethod
from PIL.ImageFile import ImageFile
import numpy.typing as npt
from PIL import Image
from frame_generator.image_utils import get_formatted_image_data


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


class TwoPlaneSolvro(Animator):
    def internal_setup(self, cube_width: int) -> None:
        img: ImageFile = Image.open('frame_generator/solvro8x8x8.png')

        self.data: list[list[float]] = get_formatted_image_data(img)
        self.counter = 0
        self.plane_index: int = cube_width // 2

    def start_function(self, leds):
        return self.update_function(leds)

    def update_function(self, leds):
        self.counter += 1
        is_even_frame = self.counter % 2 == 0

        for x in range(leds.shape[0]):
            for y in range(leds.shape[1]):
                for z in range(leds.shape[2]):
                    leds[x][y][z] = [0, 0, 0]

                    if is_even_frame and z == self.plane_index:
                        leds[x][y][z] = self.data[x + y * leds.shape[0]]
                    elif not is_even_frame and x == self.plane_index:
                        leds[x][y][z] = self.data[z + y * leds.shape[0]]

        return leds

    def stop_function(self, leds, frames):
        return self.counter > 2


class SimpleRotatingSolvro(Animator):
    def internal_setup(self, cube_width: int) -> None:
        img: ImageFile = Image.open('frame_generator/solvro8x8x8.png')
        self.data: list[list[float]] = get_formatted_image_data(img)
        self.counter = 0
        self.middle = cube_width // 2

    def start_function(self, leds):
        return self.update_function(leds)

    def update_function(self, leds):
        self.counter += 1
        mode = self.counter % 4

        for x in range(leds.shape[0]):
            for y in range(leds.shape[1]):
                for z in range(leds.shape[2]):
                    leds[x][y][z] = [0, 0, 0]

        width = leds.shape[0]

        if mode == 0:
            for x in range(width):
                for y in range(width):
                    index = x + y * width
                    if index < len(self.data):
                        leds[x][y][self.middle] = self.data[index]
        elif mode == 1:
            for z in range(width):
                for y in range(width):
                    mapped_x = width - 1 - z
                    index = mapped_x + y * width
                    if index < len(self.data):
                        leds[self.middle][y][z] = self.data[index]
        elif mode == 2:
            for x in range(width):
                for y in range(width):
                    mapped_x = width - 1 - x
                    index = mapped_x + y * width
                    if index < len(self.data):
                        leds[x][y][self.middle] = self.data[index]
        elif mode == 3:
            for z in range(width):
                for y in range(width):
                    index = z + y * width
                    if index < len(self.data):
                        leds[self.middle][y][z] = self.data[index]

        return leds

    def stop_function(self, leds, frames):
        return self.counter > 16
