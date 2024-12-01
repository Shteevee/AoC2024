defmodule Day1 do
  def stringListToIntList(list) do
    Enum.map(list, &(Integer.parse(&1) |> elem(0)))
  end

  def parseFile(file) do
    String.split(file, "\n") |> Enum.map(&(String.split(&1)) |> stringListToIntList)
  end

  def pairSmallestValues(list) do
    first = Enum.map(list, fn [x, _] -> x end) |> Enum.sort
    second = Enum.map(list, fn [_, x] -> x end) |> Enum.sort
    Enum.zip(first, second)
  end

  def sumPairDistances(list) do
    Enum.reduce(list, 0, fn {x, y}, acc -> acc + abs(x-y) end)
  end

  def createOccurrenceMap(list) do
    Enum.reduce(list, Map.new(), fn [_, x], acc -> Map.put(acc, x, Map.get(acc, x, 0)+1)  end)
  end

  def calcSimilarityScore(occurrenceMap, list) do
    Enum.reduce(list, 0, fn [x, _], acc -> acc + (x * Map.get(occurrenceMap, x, 0)) end)
  end
end

input = File.read!("input.txt") |> Day1.parseFile
summedSmallestPairs = Day1.pairSmallestValues(input) |> Day1.sumPairDistances
IO.inspect(summedSmallestPairs, label: "Part 1")
similarityScore = Day1.createOccurrenceMap(input) |> Day1.calcSimilarityScore(input)
IO.inspect(similarityScore, label: "Part 2")
