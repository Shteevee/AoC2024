defmodule Day11 do
  def string_list_to_int_list(list) do
    Enum.map(list, &(Integer.parse(&1) |> elem(0)))
  end

  def parse_file(file) do
    file
    |> String.split(" ")
    |> string_list_to_int_list()
    |> Enum.reduce(
      Map.new(),
      fn stone, map ->
        Map.put(map, stone, 1 + Map.get(map, stone, 0))
      end
    )
  end

  def update_map_split_key(map, key_string, key_length, value) do
    half_len = round(key_length / 2)

    {first, _} =
      key_string
      |> String.slice(0, half_len)
      |> Integer.parse()

    {second, _} =
      key_string
      |> String.slice(half_len, half_len)
      |> Integer.parse()

    map
    |> Map.update(first, value, fn x -> value + x end)
    |> Map.update(second, value, fn y -> value + y end)
  end

  def blink(stone_map) do
    stone_map
    |> Enum.reduce(
      Map.new(),
      fn {k, v}, map ->
        key_string = Integer.to_string(k)
        key_length = String.length(key_string)

        cond do
          k == 0 ->
            Map.put(map, 1, v + Map.get(map, 1, 0))

          rem(key_length, 2) == 0 ->
            map
            |> update_map_split_key(key_string, key_length, v)

          true ->
            Map.put(map, k * 2024, v + Map.get(map, k * 2024, 0))
        end
      end
    )
  end

  def count_stones_after_blinks(stone_map, blinks) do
    1..blinks
    |> Enum.reduce(
      stone_map,
      fn _, map ->
        blink(map)
      end
    )
    |> Enum.map(fn {_, v} -> v end)
    |> Enum.sum()
  end
end

input = File.read!("input.txt") |> Day11.parse_file()
IO.puts("Part 1: #{Day11.count_stones_after_blinks(input, 25)}")
IO.puts("Part 2: #{Day11.count_stones_after_blinks(input, 75)}")
