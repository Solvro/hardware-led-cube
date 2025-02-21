import os


class Saver:
    def __init__(self, animation_name: str) -> None:
        self.animation_name: str = animation_name

    def save(self, parsed_content: str, extension: str, path: str) -> None:
        file_name: str = f"{self.animation_name}.{extension}"

        full_path: str = os.path.join(path, file_name)

        with open(full_path, "w") as file:
            file.write(parsed_content)
