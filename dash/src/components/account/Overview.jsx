import React, { useEffect } from "react";
import { useTitleContext } from "../../state/title-context";
import Watchlists from "./Watchlists";
import EconomicCalendar from "./EconomicCalendar";

const Overview = () => {
  const { setTitle } = useTitleContext();
  useEffect(() => setTitle("Overview"), [setTitle]);
  return (
    <div>
      <h2>Overview</h2>
      <Watchlists />
      <EconomicCalendar />
    </div>
  );
};

export default Overview;
