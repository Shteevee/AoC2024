defmodule Day6 do
  def build_floor_map(line, y, floor_map) do
    Enum.reduce(
      Enum.with_index(line),
      floor_map,
      fn {char, x}, map -> Map.put(map, {x, y}, char) end
    )
  end

  def find_guard_start(floor_map) do
    Enum.find_value(floor_map, fn {k, v} -> if v == "^", do: k end)
  end

  def parse_file(file) do
    lines = String.split(file, "\n") |> Enum.map(&String.graphemes(&1))

    floor_map =
      Enum.reduce(
        Enum.with_index(lines),
        Map.new(),
        fn {line, y}, map -> build_floor_map(line, y, map) end
      )

    {find_guard_start(floor_map), floor_map}
  end

  def next_direction(current_direction) do
    case current_direction do
      {0, -1} -> {1, 0}
      {1, 0} -> {0, 1}
      {0, 1} -> {-1, 0}
      {-1, 0} -> {0, -1}
    end
  end

  def walk({{x, y}, {move_x, move_y}}, floor_map, walked \\ MapSet.new()) do
    cond do
      Map.get(floor_map, {x, y}) == nil ->
        {walked, false}

      MapSet.member?(walked, {{x, y}, {move_x, move_y}}) ->
        {walked, true}

      true ->
        new_walked = MapSet.put(walked, {{x, y}, {move_x, move_y}})
        next_pos = {x + move_x, y + move_y}

        if Map.get(floor_map, next_pos) == "#" do
          walk(
            {{x, y}, next_direction({move_x, move_y})},
            floor_map,
            new_walked
          )
        else
          walk({next_pos, {move_x, move_y}}, floor_map, new_walked)
        end
    end
  end

  def find_tiles_walked({guard_start, floor_map}) do
    walk({guard_start, {0, -1}}, floor_map)
  end

  def count_unique_tiles_walked(input) do
    {tiles_walked, _} = Day6.find_tiles_walked(input)

    Enum.reduce(tiles_walked, MapSet.new(), fn {pos, _}, set -> MapSet.put(set, pos) end)
    |> MapSet.size()
  end

  def count_object_loop_positions({guard_start, floor_map}) do
    Enum.count(floor_map, fn {pos, char} ->
      if char != "#" do
        {_, loops} = walk({guard_start, {0, -1}}, Map.put(floor_map, pos, "#"))
        loops
      end
    end)
  end
end

input = File.read!("input.txt") |> Day6.parse_file()
IO.puts("Part 1: #{Day6.count_unique_tiles_walked(input)}")
IO.puts("Part 2: #{Day6.count_object_loop_positions(input)}")
