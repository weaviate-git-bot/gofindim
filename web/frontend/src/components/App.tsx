import { BrowserRouter, HashRouter, Route, Routes } from "react-router-dom";
import Browse from "./browse";
import Main from "./main";
import Similar from "./similar";
import Header from "./header";
import { Fragment } from "react";
const App = () => {
  return (
    <HashRouter>
      <Header tw='mb-5' />
      <Routes>
        <Route path="/">
          <Route path="" element={<Main />} />
          <Route path="similar" element={<Similar />} />
          <Route path="browse" element={<Browse />} />
        </Route>
      </Routes>
    </HashRouter>
  );
};

export default App;
