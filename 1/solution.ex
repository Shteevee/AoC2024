defmodule Day1 do
  def string_list_to_int_list(list) do
    Enum.map(list, &(Integer.parse(&1) |> elem(0)))
  end

  def parse_file(file) do
    String.split(file, "\n") |> Enum.map(&(String.split(&1) |> string_list_to_int_list))
  end

  def pair_smallest_values(list) do
    first = Enum.map(list, fn [x, _] -> x end) |> Enum.sort()
    second = Enum.map(list, fn [_, x] -> x end) |> Enum.sort()
    Enum.zip(first, second)
  end

  def sum_pair_distances(list) do
    Enum.reduce(list, 0, fn {x, y}, acc -> acc + abs(x - y) end)
  end

  def create_occurrence_map(list) do
    Enum.reduce(list, Map.new(), fn [_, x], acc -> Map.put(acc, x, Map.get(acc, x, 0) + 1) end)
  end

  def calc_similarity_score(occurrence_map, list) do
    Enum.reduce(list, 0, fn [x, _], acc -> acc + x * Map.get(occurrence_map, x, 0) end)
  end
end

input = File.read!("input.txt") |> Day1.parse_file()
summed_smallest_pairs = Day1.pair_smallest_values(input) |> Day1.sum_pair_distances()
IO.puts("Part 1: #{summed_smallest_pairs}")
similarity_score = Day1.create_occurrence_map(input) |> Day1.calc_similarity_score(input)
IO.puts("Part 2: #{similarity_score}")
