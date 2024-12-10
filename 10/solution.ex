defmodule Day10 do
  def insert_to_map(map, pos, char) do
    {height, _} = Integer.parse(char)
    Map.put(map, pos, height)
  end

  def parse_file(file) do
    file
    |> String.split("\n")
    |> Enum.map(fn line -> String.graphemes(line) end)
    |> Enum.with_index()
    |> Enum.reduce(
      Map.new(),
      fn {line, y}, map ->
        line
        |> Enum.with_index()
        |> Enum.reduce(
          map,
          fn {char, x}, acc ->
            insert_to_map(acc, {x, y}, char)
          end
        )
      end
    )
  end

  def find_starting_points(height_map) do
    height_map
    |> Enum.filter(fn {_, h} -> h == 0 end)
    |> Enum.map(fn {p, _} -> p end)
  end

  def find_valid_next_pos({x, y}, height_map) do
    [{x, y + 1}, {x, y - 1}, {x - 1, y}, {x + 1, y}]
    |> Enum.filter(
        fn next_pos ->
          cond do
            Map.get(height_map, next_pos) == nil ->
              false
            Map.get(height_map, next_pos) - Map.get(height_map, {x, y}) == 1 ->
              true
            true ->
              false
          end
        end
      )
  end

  def trailheads_from_point({x, y}, height_map, endpoints \\ MapSet.new()) do
    current_height = Map.get(height_map, {x, y})
    if current_height == 9 do
      MapSet.put(endpoints, {x, y})
    else
      {x, y}
      |> find_valid_next_pos(height_map)
      |> Enum.reduce(
        endpoints,
        fn next_pos, endpoint_set ->
          MapSet.union(trailheads_from_point(next_pos, height_map, endpoint_set), endpoint_set)
        end
      )
    end
  end

  def sum_set_sizes(sets) do
    sets
    |> Enum.reduce(
      0,
      fn set, acc ->
        acc + MapSet.size(set)
      end
    )
  end

  def collect_trailhead_sets(starting_points, height_map) do
    starting_points
    |> Enum.reduce(
      [],
      fn start, list ->
        [trailheads_from_point(start, height_map) | list]
      end
    )
  end

  def count_trailheads(height_map) do
    height_map
    |> find_starting_points()
    |> collect_trailhead_sets(height_map)
    |> sum_set_sizes()
  end

  def trails_from_point({x, y}, height_map, total \\ 0) do
    current_height = Map.get(height_map, {x, y})
    if current_height == 9 do
      total + 1
    else
      {x, y}
      |> find_valid_next_pos(height_map)
      |> Enum.reduce(
          total,
        fn next_pos, acc ->
          trails_from_point(next_pos, height_map, acc)
        end
      )
    end
  end

  def count_trails(height_map) do
    height_map
    |> find_starting_points()
    |> Enum.map(
      fn start ->
        trails_from_point(start, height_map)
      end
    )
    |> Enum.sum()
  end
end

input = File.read!("input.txt") |> Day10.parse_file()
IO.puts("Part 1: #{Day10.count_trailheads(input)}")
IO.puts("Part 2: #{Day10.count_trails(input)}")
