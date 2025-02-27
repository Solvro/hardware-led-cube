from frame_generator.state_parser import PrototypeJsonifier
from frame_generator.animator import *
from frame_generator.parsed_handler import Saver
from frame_generator.main import generate_frames

ANIMATION_NAME = "animation"

if __name__ == "__main__":
    width = 8
    state_parser = PrototypeJsonifier()
    sample_animator = StaticSolvro(width)
    saver = Saver(ANIMATION_NAME)
    generate_frames((width, width, width), state_parser,
                    sample_animator, saver)
