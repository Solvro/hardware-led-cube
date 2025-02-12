from frame_generator.state_parser import PrototypeJsonifier
from frame_generator.animator import *
from frame_generator.parsed_handler import Saver
from frame_generator.main import generate_frames

if __name__ == "__main__":
    state_parser = PrototypeJsonifier()
    sample_animator = SampleAnimator()
    saver = Saver()
    generate_frames(4, 4, 4, state_parser, sample_animator, saver)
