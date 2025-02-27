
def get_formatted_image_data(img) -> list[list[float]]:
    data = list(img.getdata())
    data.reverse()
    data = [[c/255 for c in color[:3]] for color in data]
    return data
