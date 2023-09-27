import Similar from "components/similar";
import {Link} from "react-router-dom";

const Main = () => {

  return (
    <div tw='place-content-center items-center justify-around w-full flex flex-row'>
    <Link to="similar">
    View Similar
    </Link>
    <Link to="browse">
    Browse
    </Link>
    </div>
  );
};

export default Main;
