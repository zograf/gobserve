import React, { useState } from "react";
import Chart from "react-apexcharts";

export default function PieChartCard({title, x, y, width, colors, height}) {
    const [data, setData] = useState({
      options: {
        chart: {
          id: "basic-line",
        },
        labels: y,
        colors: colors
      },
        series: x
    })

    return (
        <div className="card" style={{width: width, height: height}}>
            <p className="chart-title">{title}</p>
            <Chart
                options={data.options}
                series={data.series}
                type="pie"
                height={"100%"}
                width={"100%"}
            />
        </div>
    );
}