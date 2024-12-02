defmodule Day2 do
  def stringListToIntList(list) do
    Enum.map(list, &(Integer.parse(&1) |> elem(0)))
  end

  def parseFile(file) do
    String.split(file, "\n") |> Enum.map(&(String.split(&1)) |> stringListToIntList)
  end

  def diffInBounds(x, y) do
    diff = x - y
    diff >= 1 && diff <= 3
  end

  def inBounds([x, y | tail], func) do
    inBounds([y | tail], func, func.(x, y))
  end

  def inBounds([x, y | tail], func, acc) do
    newAcc = func.(x, y) && acc
    if length(tail) > 0 do
      inBounds([y | tail], func, newAcc)
    else
      newAcc
    end
  end

  def reportIsSafe(report) do
    diffInBoundsDsc = fn x, y -> diffInBounds(x, y) end
    diffInBoundsAsc = fn x, y -> diffInBounds(y, x) end
    inBounds(report, diffInBoundsAsc) || inBounds(report, diffInBoundsDsc)
  end

  def countSafeReports(reports) do
    Enum.reduce(
      reports,
      0,
      fn report, acc -> if reportIsSafe(report) do acc + 1 else acc end end
    )
  end

  def reportIsSafeWithTolerance(report) do
    reportVariants = [report |
      Enum.reduce(
        0..length(report)-1,
        [],
        fn idx, acc -> [ List.delete_at(report, idx) | acc ] end
      )
    ]

    Enum.map(reportVariants, &(reportIsSafe(&1))) |> Enum.any?
  end

  def countSafeReportsWithTolerance(reports) do
    Enum.reduce(
      reports,
      0,
      fn report, acc -> if reportIsSafeWithTolerance(report) do acc + 1 else acc end end
    )
  end
end

reports = File.read!("input.txt") |> Day2.parseFile
IO.inspect(Day2.countSafeReports(reports), label: "Part 1")
IO.inspect(Day2.countSafeReportsWithTolerance(reports), label: "Part 2")
