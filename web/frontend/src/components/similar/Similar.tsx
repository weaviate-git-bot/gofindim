import tw, { styled } from "twin.macro";
import { FormEvent, useEffect, useState } from "react";
import MatchedImage from "./MatchedImage";
import { useSearchParams } from "react-router-dom";

interface Image {
  name: string;
  path: string;
  id: string;
}

const roundedFloatFromString = (num: string, precision = 2) =>
  Math.round(parseFloat(num) * Math.pow(10, precision)) /
  Math.pow(10, precision);

const Similar = () => {
  const [images, setImages] = useState([] as Image[]);
  const [imagesToDelete, setImagesToDelete] = useState([] as string[]);
  const [searchParams, setSearchParams] = useSearchParams();
  const [settings, setSettings] = useState({
    limit: 10,
    distance: 0.7,
    image_weight: 0.5,
    text_weight: 0.5,
    text_input: searchParams.get("text_input") || "",
    path: "",
    uuid: searchParams.get("uuid") || "",
  });

  const clearCheckboxes = () => {
    Array.from(document.getElementsByClassName("delete_images_check")).forEach(
      (e) => {
        if (e instanceof HTMLInputElement) {
          e.checked = false;
        }
      }
    );
    setImagesToDelete([]);
  };

  useEffect(() => {
    let url = `${config.BaseURL}/api/similar?&path=${searchParams.get( "path"
    )}&limit=${settings.limit
    }&distance=${settings.distance
    }&image_weight=${ settings.image_weight 
    }&text_weight=${settings.text_weight
    }&text_input=${settings.text_input}`;
    if (searchParams.get("uuid") != null) {
      url += `&uuid=${searchParams.get("uuid")}`;
    }
    if (searchParams.get("path") || searchParams.get("uuid")) {
      fetch(url).then((res) => {
        res
          .json()
          .then((data) => {
            setImages(data.images || []);
          })
          .catch((err) => {
            console.log(err);
          });
      });
    }
  }, [searchParams]);

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const formData = new FormData(e.target as HTMLFormElement);
    imagesToDelete.forEach((img) => {
      formData.append("delete_images[]", img);
    });
    Object.keys(settings).forEach((key) => {
      formData.set(key, settings[key as keyof typeof settings].toString());
    });

    fetch("http://localhost:8888/api/similar", {
      method: "POST",
      body: formData,
    }).then((res) => {
      res.json().then((data) => {
        console.log(data);
        setImages(data.images || []);
      });
    });
    setSearchParams({
      text_input: formData.get("text_input") as string,
      limit: settings.limit.toString(),
      distance: settings.distance.toString(),
      image_weight: settings.image_weight.toString(),
      text_weight: settings.text_weight.toString(),
      path: searchParams.get("path") || "",
    });
    clearCheckboxes();
  };

  const handleCheck = (e: React.ChangeEvent<HTMLInputElement>) => {
    setImagesToDelete((prev) => {
      if (e.target.checked) {
        return [...prev, e.target.value];
      } else {
        return prev.filter((uuid) => uuid !== e.target.value);
      }
    });
  };

  const handleApplySettings = () => {
    const limit = settings.limit.toString();
    const distance = settings.distance.toString();
    const text_input = settings.text_input.toString();
    const image_weight = settings.image_weight.toString();
    const text_weight = settings.text_weight.toString();
    const path = searchParams.get("path") || "";
    const uuid = searchParams.get("uuid") || "";
    const newParams: Record<string, any> = {
      path,
      limit,
      distance,
      text_input,
      image_weight,
      text_weight,
    };
    if (uuid) {
      newParams.uuid = uuid;
    }
    setSearchParams((prev) => {
      return newParams;
    });
  };

  return (
    <div tw="place-content-center items-center w-full">
      <div tw="items-center justify-center w-full font-['Inter'] relative">
        {imagesToDelete.map((uuid) => (
          <div key={uuid}>{uuid}</div>
        ))}
        <div tw="grid w-full h-auto mt-[15em] grid-flow-row  grid-cols-2 md:grid-cols-3 xl:grid-cols-6 items-end ">
          {images.map((image, i) => (
            <MatchedImage
              key={i}
              src={`${image.path}`}
              alt={`similar-${i}`}
              uuid={image.id}
              handleCheck={handleCheck}
              onClick={() => {
                setSearchParams({
                  uuid: image.id,
                  limit: settings.limit.toString(),
                  distance: settings.distance.toString(),
                });
              }}
            />
          ))}
        </div>
        <div tw="fixed top-10 left-[50%] translate-x-[-50%] w-full">
          <form
            tw="items-center place-content-center flex flex-row ring-2 ring-slate-200 rounded-sm m-2 bg-[#333] "
            onSubmit={(e) => handleSubmit(e)}
          >
            <InputGrid>
              <div tw="flex flex-col gap-y-2">
                <Field>
                  <label htmlFor="text_input">Text:</label>
                  <Input
                    type="text"
                    name="text_input"
                    onChange={(e) =>
                      setSettings((prev) => {
                        return {
                          ...prev,
                          text_input: e.target.value,
                        };
                      })
                    }
                    value={settings.text_input}
                  />
                </Field>
                <Field>
                  <label htmlFor="file_input">File:</label>
                  <Input type="file" name="image_input" />
                </Field>
              </div>
              <div tw="flex flex-col gap-y-2">
                <Field>
                  <label htmlFor="text_weight">
                    Text Weight: {settings.text_weight.toFixed(2)}
                  </label>
                  <Range
                    tw="w-[10em]"
                    type="range"
                    step={0.01}
                    max={1}
                    min={0}
                    value={settings.text_weight}
                    name="text_weight"
                    onChange={(e) =>
                      setSettings((prev) => {
                        return {
                          ...prev,
                          text_weight: roundedFloatFromString(e.target.value),
                          image_weight:
                            Math.round((1 - parseFloat(e.target.value)) * 100) /
                            100,
                        };
                      })
                    }
                  />
                </Field>
                <Field>
                  <label htmlFor="image_weight">
                    Image Weight:{settings.image_weight.toFixed(2)}
                  </label>
                  <Range
                    tw="w-[10em]"
                    type="range"
                    step={0.01}
                    max={1}
                    min={0}
                    name="image_weight"
                    value={settings.image_weight}
                    onChange={(e) =>
                      setSettings((prev) => {
                        return {
                          ...prev,
                          image_weight: roundedFloatFromString(e.target.value),
                          text_weight:
                            Math.round((1 - parseFloat(e.target.value)) * 100) /
                            100,
                        };
                      })
                    }
                  />
                </Field>
              </div>
              <div tw="flex flex-col gap-y-2">
                <Field>
                  <label htmlFor="distance">
                    Distance: {settings.distance}
                  </label>
                  <Range
                    min={0}
                    max={1}
                    defaultValue={0.8}
                    step={0.01}
                    type="range"
                    name="distance"
                    onChange={(e) => {
                      setSettings((prev) => {
                        return {
                          ...prev,
                          distance: parseFloat(e.target.value),
                        };
                      });
                    }}
                  />
                </Field>

                <Field>
                  <label htmlFor="limit">Limit: {settings.limit}</label>
                  <Range
                    defaultValue={10}
                    step={1}
                    type="range"
                    min={1}
                    max={50}
                    name="limit"
                    onChange={(e) => {
                      setSettings((prev) => {
                        return { ...prev, limit: parseInt(e.target.value) };
                      });
                    }}
                  />
                </Field>
              </div>

              <Field tw="row-start-2 col-span-4 place-self-center">
                <div tw="flex flex-row">
                  <Button
                    type="button"
                    onClick={() =>
                      setSettings((prev) => {
                        return {
                          ...prev,
                          image_weight: 0.5,
                          text_weight: 0.5,
                        };
                      })
                    }
                  >
                    Reset
                  </Button>
                  <SubmitButton type="submit" value="Submit" />
                  <Button type="button" onClick={() => handleApplySettings()}>
                    Apply
                  </Button>
                </div>
              </Field>
            </InputGrid>
            <div tw="flex flex-row gap-x-2"></div>
          </form>
        </div>
      </div>
    </div>
  );
};

const Button = styled("button")(() => [
  tw`bg-slate-200 text-slate-800 rounded-sm p-1 m-1 cursor-pointer`,
]);

const SubmitButton = styled("input")(() => [
  tw`bg-slate-200 text-slate-800 rounded-sm p-1 m-1 cursor-pointer`,
]);

const InputGrid = styled("div")(() => [
  tw`grid grid-cols-[5fr 4fr 3fr 1fr ] m-2 p-2 w-full gap-x-2 place-content-around`,
]);

const Range = styled("input")(() => [tw`w-[7em]`]);

const Input = styled("input")(() => [tw`w-full`]);

const Field = styled("div")(() => [tw`flex flex-row justify-between`]);

const NumberInput = styled(Input)(() => [tw``]);

export default Similar;
