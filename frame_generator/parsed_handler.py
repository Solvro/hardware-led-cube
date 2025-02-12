import os


class Saver:
    def save(self, parsed_content: str, frame: int, animation_name: str, path: str) -> None:
        file_name: str = f"{animation_name}_frame_{frame}.json"

        full_path = os.path.join(path, file_name)

        with open(full_path, "+a") as file:
            file.write(parsed_content)
