import React, { useState, useEffect, useRef } from "react";
import Chart from "react-apexcharts";

export default function PieChartCard({title, x, y, width, colors, height}) {
    const [data, setData] = useState({
        options: {
          chart: {
            id: "basic-pie",
          },
          labels: y,
          colors: colors
        },
          series: x
    })

    useEffect(() => {
      setData(prevData => ({
        options: {
          chart: prevData.options.chart,
          labels: y,
          colors: prevData.options.colors
        },
        series: x
    }));
    }, [x, y])

    return (
        <div className="card" style={{width: width, height: height}}>
            <p className="chart-title">{title}</p>
            <Chart
                options={data.options}
                series={data.series}
                type="pie"
                height={"95%"}
                width={"100%"}
            />
        </div>
    );
}