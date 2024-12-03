defmodule Day3 do
  def mult_match_to_mult_pair([x, y]) do
    {int_x, _} = Integer.parse(x)
    {int_y, _} = Integer.parse(y)
    {int_x, int_y}
  end

  def parse_mult_pairs(string) do
    regex = ~r"mul\((\d{1,3})\,(\d{1,3})\)"

    Regex.scan(regex, string)
    |> Enum.map(fn [_ | pair] -> mult_match_to_mult_pair(pair) end)
  end

  def sum_mult_pairs(pairs) do
    Enum.reduce(pairs, 0, fn {x, y}, acc -> acc + x * y end)
  end

  def grab_do_blocks([head | tail]) do
    Enum.reduce(
      tail,
      head,
      fn instr_block, acc ->
        case instr_block do
          [_] ->
            acc

          [_ | do_block] ->
            acc ++ do_block
        end
      end
    )
  end

  def find_enabled_blocks(file) do
    String.split(file, "don't()") |> Enum.map(&String.split(&1, "do()")) |> grab_do_blocks()
  end

  def sum_enabled_mult_pairs(file) do
    find_enabled_blocks(file)
    |> Enum.reduce([], fn str, acc -> Enum.concat(acc, parse_mult_pairs(str)) end)
    |> sum_mult_pairs()
  end
end

file = File.read!("input.txt")
IO.puts("Part 1: #{Day3.parse_mult_pairs(file) |> Day3.sum_mult_pairs()}")
IO.puts("Part 2: #{Day3.sum_enabled_mult_pairs(file)}")
