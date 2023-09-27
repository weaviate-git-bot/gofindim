import { Fragment, FunctionComponent, useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import PageNavigator from "./PageNavigator";
import { PrimaryButton } from "components/Buttons";

interface File {
  name: string;
  isDir: boolean;
  thumbnailUrl: string;
  id: string;
  path: string;
}
function PathButton({ path, onClick }: { path: string; onClick: () => void }) {
  return <button onClick={onClick}>{path}</button>;
}

const Browse = () => {
  const [files, setFiles] = useState<File[]>([]);
  let [searchParams, setSearchParams] = useSearchParams({
    path: "/",
    page: "0",
    limit: "15",
  });
  const [limit, setLimit] = useState(15);
  const [pageCount, setPageCount] = useState(0);

  function PathNavigator({
    fullPath,
  }: {
    fullPath: string;
  }) {
    const parts = fullPath.split("/").filter((part) => part !== "");

    return (
      <div>
        {parts.map((part: string, idx: number) => {
          const path = "/" + parts.slice(0, idx + 1).join("/");
          return (
            <Fragment key={path}>
              {idx > 0 && <span>/</span>}
              <PathButton
                path={part}
                onClick={() => {
                  setPath(path);
                }}
              />
            </Fragment>
          );
        })}
      </div>
    );
  }

  const mkFileUrlPath = (file: File) => {
    const path = (searchParams.get("path") || "") + "/" + file.name;
    console.log(path);
    setSearchParams({ path: path, page: "0", limit: limit.toString() });
  };

  const setPath = (path: string) => {
    setSearchParams({ path: path, page: "0", limit: limit.toString() });
  };

  const scanHandler = (file: File) => {
    fetch(config.BaseURL + "/api/scan?path=" + encodeURIComponent(file.path), {
      method: "POST",
    })
      .then(() => {
        console.log("OK");
      })
      .catch((err) => {
        console.log("NOT OK", err);
      });
  };
  const goUpHandler = () => {
    const inputPath = searchParams.get("path") || "";
    const lastSlashIndex = inputPath.lastIndexOf("/");
    const truncatedPath = inputPath.slice(0, lastSlashIndex);
    setSearchParams({ path: truncatedPath });
  };

  const setPage = (val: number) => {
    setSearchParams({
      path: searchParams.get("path") || "",
      page: val.toString(),
      limit: limit.toString(),
    });
  };

  useEffect(() => {
    setSearchParams({
      path: searchParams.get("path") || "",
      page: searchParams.get("page") || "0",
      limit: limit.toString(),
    });
  }, [pageCount, limit, limit]);

  useEffect(() => {
    console.log(config.BaseURL);
    const queryStrings = searchParams.toString();
    fetch(config.BaseURL + "/api/files?" + queryStrings).then((res) => {
      res
        .json()
        .then((data) => {
          const nonNullFiles = data.files.filter((file: File) => file != null);
          setFiles(nonNullFiles || []);
          setPageCount(data.itemCount / limit);
        })
        .catch((err) => {
          setFiles([]);
          console.log(err);
        });
    });
  }, [searchParams]);
  return (
    <div className="browse">
      <PathNavigator
        fullPath={searchParams.get("path") || ""}
      />
      <div tw="grid grid-cols-3 w-full object-center">
        <button
          tw="bg-slate-500 rounded-md p-2 flex w-fit place-self-center align-middle mt-2"
          onClick={goUpHandler}
        >
          Go Up
        </button>
        <PageNavigator
          tw="place-self-center items-center"
          page={parseInt(searchParams.get("page") || "0")}
          setPage={setPage}
          pageCount={pageCount}
          pageSize={limit}
        />
        <div tw="flex flex-col items-center place-content-center col-start-3">
          <label htmlFor="limit">Limit: {limit}</label>
          <input
            defaultValue={10}
            step={1}
            type="range"
            min={1}
            max={50}
            name="limit"
            onChange={(e) => {
              setLimit((prev) => {
                return parseInt(e.target.value);
              });
            }}
          />
        </div>
      </div>
      <div tw="mt-[3em] grid grid-cols-5 gap-4 place-items-center items-start">
        {files.map((file, idx) => (
          <div
            tw="flex flex-col items-center place-content-around w-fit h-full"
            key={idx}
          >
            {file.thumbnailUrl && (
              <Link target="_blank" to={config.BaseURL+file.path}>
                <img
                  src={config.BaseURL+file.thumbnailUrl}
                  alt={file.name}
                />
              </Link>
            )}
            {!file.thumbnailUrl && (
              <PlaceholderSVG
                onClick={() => {
                  mkFileUrlPath(file);
                }}
              />
            )}
            <button
              onClick={() => {
                mkFileUrlPath(file);
              }}
            >
              <p tw="w-fit max-w-[150px] text-xs break-words">{file.name}</p>
            </button>
            <div tw="flex flex-row gap-x-2">
              <Link to={"/similar?path=" + encodeURIComponent(file.path)}>
                <PrimaryButton >Similar</PrimaryButton>
              </Link>
              <PrimaryButton
                onClick={() => {
                  scanHandler(file);
                }}
              >
                Scan
              </PrimaryButton>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

const PlaceholderSVG: FunctionComponent<Record<string, any>> = (props) => {
  return (
    <svg
      tw="animate-[pulse_2s] cursor-pointer"
      xmlns="http://www.w3.org/2000/svg"
      width="100"
      height="100"
      fill="tan"
      viewBox="0 0 24 24"
      {...props}
    >
      <path d="M10 4H4a2 2 0 00-2 2v12a2 2 0 002 2h16a2 2 0 002-2V8a2 2 0 00-2-2h-8l-2-2z" />
    </svg>
  );
  // return (
  //   <svg

  //     tw="animate-pulse"
  //     width="100px"
  //     height="100px"
  //     viewBox="0 0 100 100"
  //     preserveAspectRatio="none"
  //   >
  //     <rect x="0" y="0" width="100" height="100" />
  //   </svg>
  // );
};

export default Browse;
