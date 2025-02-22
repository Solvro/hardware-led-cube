import numpy as np
import numpy.typing as npt
from frame_generator.state_parser import StateParser
from frame_generator.animator import *
from frame_generator.parsed_handler import Saver
from typing import Tuple

COLOR_CHANNELS_LENGTH = 3
ANIMATION_DIRECTORY = ""


class LedCube:
    def __init__(self, width, height, depth, state_parser: StateParser) -> None:
        self.leds: npt.NDArray = np.zeros(
            (width, height, depth, COLOR_CHANNELS_LENGTH))
        self.state_parser: StateParser = state_parser

        self.states: list[npt.NDArray] = []
        self.frame: int = 0

    def start(self, start_function) -> None:
        self.leds = start_function(self.leds)
        self.states = [self.leds]

    def update(self, update_function) -> None:
        update_function(self.leds)
        self.states += [self.leds]

    def update_frame(self, frame: int) -> None:
        self.frame = frame

    def should_stop(self, stop_function) -> bool:
        return stop_function(self.leds, self.frame)

    def parse(self) -> None:
        self.state_parser.handle_frame(self.leds)


def generate_frames(dimensions: Tuple[int, int, int], state_parser: StateParser, animator: Animator, saver: Saver) -> None:
    led_cube = LedCube(*dimensions, state_parser)

    led_cube.start(animator.start_function)

    frame = 0

    def is_running() -> bool:
        if animator.stop_function is None:
            return True

        return not led_cube.should_stop(animator.stop_function)

    while is_running():
        led_cube.parse()
        print(f"Generated frame {frame}")

        led_cube.update(animator.update_function)
        led_cube.update_frame(frame)
        frame += 1

    saver.save(state_parser.get_parsed_results(),
               state_parser.extension, ANIMATION_DIRECTORY)
