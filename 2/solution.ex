defmodule Day2 do
  def string_list_to_int_list(list) do
    Enum.map(list, &(Integer.parse(&1) |> elem(0)))
  end

  def parse_file(file) do
    String.split(file, "\n") |> Enum.map(&(String.split(&1) |> string_list_to_int_list))
  end

  def diff_in_bounds(x, y) do
    diff = x - y
    diff >= 1 && diff <= 3
  end

  def in_bounds([x, y | tail], func) do
    in_bounds([y | tail], func, func.(x, y))
  end

  def in_bounds([x, y | tail], func, acc) do
    new_acc = func.(x, y) && acc

    if length(tail) > 0 do
      in_bounds([y | tail], func, new_acc)
    else
      new_acc
    end
  end

  def report_is_safe(report) do
    diff_in_bounds_dsc = fn x, y -> diff_in_bounds(x, y) end
    diff_in_bounds_asc = fn x, y -> diff_in_bounds(y, x) end
    in_bounds(report, diff_in_bounds_asc) || in_bounds(report, diff_in_bounds_dsc)
  end

  def count_safe_reports(reports) do
    Enum.reduce(
      reports,
      0,
      fn report, acc ->
        if report_is_safe(report) do
          acc + 1
        else
          acc
        end
      end
    )
  end

  def report_is_safe_with_tolerance(report) do
    report_variants = [
      report
      | Enum.reduce(
          0..(length(report) - 1),
          [],
          fn idx, acc -> [List.delete_at(report, idx) | acc] end
        )
    ]

    Enum.map(report_variants, &report_is_safe(&1)) |> Enum.any?()
  end

  def count_safe_reports_with_tolerance(reports) do
    Enum.reduce(
      reports,
      0,
      fn report, acc ->
        if report_is_safe_with_tolerance(report) do
          acc + 1
        else
          acc
        end
      end
    )
  end
end

reports = File.read!("input.txt") |> Day2.parse_file()
IO.puts("Part 1: #{Day2.count_safe_reports(reports)}")
IO.puts("Part 2: #{Day2.count_safe_reports_with_tolerance(reports)}")
