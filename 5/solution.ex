defmodule Day5 do
  def string_list_to_int_list(list) do
    Enum.map(list, &(Integer.parse(&1) |> elem(0)))
  end

  def calc_center_of_list(list) do
    round((length(list) - 1) / 2)
  end

  def create_order_rule_map(rule_strings) do
    Enum.reduce(
      rule_strings,
      Map.new(),
      fn rule_string, rule_map ->
        [x, y] = String.split(rule_string, "|") |> string_list_to_int_list()
        Map.put(rule_map, x, [y | Map.get(rule_map, x, [])])
      end
    )
  end

  def create_page_list(page_strings) do
    Enum.map(page_strings, &String.split(&1, ","))
    |> Enum.map(&string_list_to_int_list(&1))
  end

  def parse_file(file) do
    line_split = String.split(file, "\n")

    order_rules =
      Enum.filter(line_split, &String.contains?(&1, "|"))
      |> create_order_rule_map()

    pages =
      Enum.filter(line_split, &String.contains?(&1, ","))
      |> create_page_list()

    {order_rules, pages}
  end

  def valid_page_list?([x | [y]], order_rules) do
    Enum.member?(Map.get(order_rules, x, []), y)
  end

  def valid_page_list?([x | [y | tail]], order_rules) do
    Enum.member?(Map.get(order_rules, x, []), y) && valid_page_list?([y | tail], order_rules)
  end

  def sum_middle_page_of_valid_lists({order_rules, page_lists}) do
    Enum.filter(page_lists, &valid_page_list?(&1, order_rules))
    |> Enum.reduce(
      0,
      fn page_list, acc ->
        acc + Enum.at(page_list, calc_center_of_list(page_list))
      end
    )
  end

  def fix_page_list([x], _, acc) do
    [x | acc]
  end

  def fix_page_list(list, order_rules, acc) do
    next_elem_idx =
      Enum.find_index(
        list,
        fn x ->
          Enum.all?(list, fn y -> x == y || Enum.member?(Map.get(order_rules, x, []), y) end)
        end
      )

    fix_page_list(
      List.delete_at(list, next_elem_idx),
      order_rules,
      [Enum.at(list, next_elem_idx) | acc]
    )
  end

  def find_fixed_page_list_center(page_list, order_rules) do
    fix_page_list(page_list, order_rules, [])
    |> Enum.at(calc_center_of_list(page_list))
  end

  def sum_middle_page_of_fixed_invalid_lists({order_rules, page_lists}) do
    Enum.filter(page_lists, &(!valid_page_list?(&1, order_rules)))
    |> Enum.reduce(
      0,
      fn page_list, acc ->
        acc + find_fixed_page_list_center(page_list, order_rules)
      end
    )
  end
end

input = File.read!("input.txt") |> Day5.parse_file()
IO.puts("Part 1: #{Day5.sum_middle_page_of_valid_lists(input)}")
IO.puts("Part 2: #{Day5.sum_middle_page_of_fixed_invalid_lists(input)}")
