defmodule Day8 do
  def insert_antena_into_map(map, char, pos) do
    if char != "." do
      Map.put(map, char, [pos | Map.get(map, char, [])])
    else
      map
    end
  end

  def parse_file(file) do
    split_file =
      file
      |> String.split("\n")
      |> Enum.map(fn line -> String.graphemes(line) end)

    antena_map =
      split_file
      |> Enum.with_index()
      |> Enum.reduce(
        Map.new(),
        fn {line, y}, map ->
          line
          |> Enum.with_index()
          |> Enum.reduce(
            map,
            fn {char, x}, acc ->
              insert_antena_into_map(acc, char, {x, y})
            end
          )
        end
      )

    {antena_map, {length(Enum.at(split_file, 0)), length(split_file)}}
  end

  def generate_antinode_positions({ax, ay}, {bx, by}, antinode_set, steps) do
    Enum.reduce(
      1..steps,
      antinode_set,
      fn step, set ->
        {x_diff, y_diff} = {ax - bx, ay - by}

        set
        |> MapSet.put({ax + step * x_diff, ay + step * y_diff})
        |> MapSet.put({bx - step * x_diff, by - step * y_diff})
      end
    )
  end

  def find_antinodes(antenas, antinodes, steps) do
    antenas
    |> Enum.reduce(
      antinodes,
      fn {ax, ay}, set ->
        antenas
        |> List.delete({ax, ay})
        |> Enum.reduce(
          set,
          fn b, acc ->
            generate_antinode_positions({ax, ay}, b, acc, steps)
          end
        )
      end
    )
  end

  def in_bounds?({x, y}, {x_bound, y_bound}) do
    x >= 0 && y >= 0 && x < x_bound && y < y_bound
  end

  def place_antinodes({antena_map, bound}, steps) do
    antena_map
    |> Enum.reduce(
      MapSet.new(),
      fn {_, antenas}, antinode_map ->
        find_antinodes(antenas, antinode_map, steps)
      end
    )
    |> MapSet.filter(fn {x, y} ->
      in_bounds?({x, y}, bound)
    end)
  end

  def find_antena_antinodes(antenas, antinode_set) do
    if length(antenas) > 1 do
      antenas
      |> Enum.reduce(
        antinode_set,
        fn antena, set ->
          MapSet.put(set, antena)
        end
      )
    end
  end

  def place_antinodes_linearly({antena_map, {x_bound, y_bound}}) do
    antinode_set = place_antinodes({antena_map, {x_bound, y_bound}}, x_bound)

    antena_map
    |> Enum.reduce(
      antinode_set,
      fn {_, antenas}, set ->
        find_antena_antinodes(antenas, set)
      end
    )
  end
end

input = File.read!("input.txt") |> Day8.parse_file()
IO.puts("Part 1: #{MapSet.size(Day8.place_antinodes(input, 1))}")
IO.puts("Part 2: #{MapSet.size(Day8.place_antinodes_linearly(input))}")
