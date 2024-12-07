defmodule Day7 do
  def string_list_to_int_list(list) do
    Enum.map(list, &(Integer.parse(&1) |> elem(0)))
  end

  def parse_file(file) do
    file
    |> String.split("\n")
    |> Enum.map(fn line -> String.split(line, ": ") end)
    |> Enum.map(fn [target, nums] ->
      {
        Integer.parse(target) |> elem(0),
        String.split(nums, " ") |> string_list_to_int_list()
      }
    end)
  end

  def achieves_target?(target, ops, nums, total) do
    case nums do
      [x | tail] ->
        ops
        |> Enum.any?(fn op -> achieves_target?(target, ops, tail, op.(total, x)) end)

      [] ->
        target == total
    end
  end

  def sum_targets_achieved(input, ops) do
    input
    |> Enum.map(fn {target, [head | tail]} ->
      if achieves_target?(target, ops, tail, head) do
        target
      else
        0
      end
    end)
    |> Enum.sum()
  end

  def sum_targets_achieved_without_concat(input) do
    ops = [
      fn x, y -> x + y end,
      fn x, y -> x * y end
    ]

    sum_targets_achieved(input, ops)
  end

  def sum_targets_achieved_with_concat(input) do
    ops = [
      fn x, y -> x + y end,
      fn x, y -> x * y end,
      fn x, y -> x * 10 ** String.length(Integer.to_string(y)) + y end
    ]

    sum_targets_achieved(input, ops)
  end
end

input = File.read!("input.txt") |> Day7.parse_file()
IO.puts("Part 1: #{Day7.sum_targets_achieved_without_concat(input)}")
IO.puts("Part 2: #{Day7.sum_targets_achieved_with_concat(input)}")
