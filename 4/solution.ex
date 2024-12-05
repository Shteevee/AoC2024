defmodule Day4 do
  def parse_file(file) do
    String.split(file, "\n")
  end

  def char_at_pos(crossword, pos) do
    {x, y} = pos
    Enum.at(crossword, y) |> String.at(x)
  end

  def find_char_positions(crossword, char) do
    find_char_positions(crossword, char, {0, 0}, [])
  end

  def find_char_positions(crossword, char, pos, acc) do
    {x, y} = pos

    new_acc =
      if char_at_pos(crossword, pos) == char do
        [pos | acc]
      else
        acc
      end

    cond do
      x < String.length(Enum.at(crossword, 0)) - 1 ->
        find_char_positions(crossword, char, {x + 1, y}, new_acc)

      y < length(crossword) - 1 ->
        find_char_positions(crossword, char, {0, y + 1}, new_acc)

      true ->
        new_acc
    end
  end

  def calc_pos_range({pos_x, pos_y}, {offset_x, offset_y}, steps) do
    Enum.map(1..steps, &{pos_x + offset_x * &1, pos_y + offset_y * &1})
  end

  # what is a matrix
  def generate_paths_to_check(crossword, x_pos, steps) do
    {x, y} = x_pos
    can_left = x - 3 >= 0
    can_right = x + 3 < String.length(Enum.at(crossword, 0))
    can_up = y - 3 >= 0
    can_down = y + 3 < length(crossword)

    paths_to_check = []

    paths_to_check =
      if can_left do
        [calc_pos_range(x_pos, {-1, 0}, steps) | paths_to_check]
      else
        paths_to_check
      end

    paths_to_check =
      if can_left && can_up do
        [calc_pos_range(x_pos, {-1, -1}, steps) | paths_to_check]
      else
        paths_to_check
      end

    paths_to_check =
      if can_up do
        [calc_pos_range(x_pos, {0, -1}, steps) | paths_to_check]
      else
        paths_to_check
      end

    paths_to_check =
      if can_up && can_right do
        [calc_pos_range(x_pos, {1, -1}, steps) | paths_to_check]
      else
        paths_to_check
      end

    paths_to_check =
      if can_right do
        [calc_pos_range(x_pos, {1, 0}, steps) | paths_to_check]
      else
        paths_to_check
      end

    paths_to_check =
      if can_right && can_down do
        [calc_pos_range(x_pos, {1, 1}, steps) | paths_to_check]
      else
        paths_to_check
      end

    paths_to_check =
      if can_down do
        [calc_pos_range(x_pos, {0, 1}, steps) | paths_to_check]
      else
        paths_to_check
      end

    paths_to_check =
      if can_down && can_left do
        [calc_pos_range(x_pos, {-1, 1}, steps) | paths_to_check]
      else
        paths_to_check
      end

    paths_to_check
  end

  def char_from_pos(crossword, {x, y}) do
    Enum.at(crossword, y) |> String.at(x)
  end

  def xmas_count(crossword, x_pos) do
    chars_remaining = 3

    generate_paths_to_check(crossword, x_pos, chars_remaining)
    |> Enum.map(fn dir_to_check ->
      Enum.reduce(dir_to_check, "", fn pos, acc -> acc <> char_from_pos(crossword, pos) end)
    end)
    |> Enum.reduce(0, fn x, acc ->
      if x == "MAS" do
        acc + 1
      else
        acc
      end
    end)
  end

  def count_xmas_occurrences(crossword) do
    find_char_positions(crossword, "X")
    |> Enum.reduce(0, &(xmas_count(crossword, &1) + &2))
  end

  def x_mas?(crossword, {x, y}) do
    can_left = x - 1 >= 0
    can_right = x + 1 < String.length(Enum.at(crossword, 0))
    can_up = y - 1 >= 0
    can_down = y + 1 < length(crossword)

    if can_left && can_right && can_up && can_down do
      top_left = char_at_pos(crossword, {x - 1, y - 1})
      bottom_right = char_at_pos(crossword, {x + 1, y + 1})
      top_right = char_at_pos(crossword, {x + 1, y - 1})
      bottom_left = char_at_pos(crossword, {x - 1, y + 1})

      ((top_left == "M" && bottom_right == "S") || (top_left == "S" && bottom_right == "M")) &&
        ((top_right == "M" && bottom_left == "S") || (top_right == "S" && bottom_left == "M"))
    end
  end

  def count_x_mas_occurrences(crossword) do
    find_char_positions(crossword, "A")
    |> Enum.reduce(
      0,
      fn pos, acc ->
        if x_mas?(crossword, pos) do
          acc + 1
        else
          acc
        end
      end
    )
  end
end

crossword = File.read!("input.txt") |> Day4.parse_file()
IO.puts("Part 1: #{Day4.count_xmas_occurrences(crossword)}")
IO.puts("Part 2: #{Day4.count_x_mas_occurrences(crossword)}")
