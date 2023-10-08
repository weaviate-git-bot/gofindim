import tw, { styled } from "twin.macro";
import { FormEvent, useCallback, useEffect, useRef, useState } from "react";
import MatchedImage from "./MatchedImage";
import { useSearchParams } from "react-router-dom";
import InputGrid from "./InputGrid/InputGrid";
import { Image } from "@/types/Image";

const Similar = () => {
  const [images, setImages] = useState({} as Record<string, Image>);
  const [imagesToDelete, setImagesToDelete] = useState([] as string[]);
  const [searchParams, setSearchParams] = useSearchParams();
  const formRef = useRef<HTMLFormElement>(null);
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
    let url = `${config.BaseURL}/api/similar?&path=${searchParams.get(
      "path"
    )}&limit=${settings.limit}&distance=${settings.distance}&image_weight=${
      settings.image_weight
    }&text_weight=${settings.text_weight}&text_input=${settings.text_input}`;
    if (searchParams.get("uuid") != null) {
      url += `&uuid=${searchParams.get("uuid")}`;
    }
    if (searchParams.get("path") || searchParams.get("uuid")) {
      fetch(url).then((res) => {
        res
          .json()
          .then((data) => {
            if (data.images) {
              const _images: Image[] = data.images;
              const rendered: Record<string, Image> = _images.reduce(
                (acc, v, idx) => {
                  acc[_images[idx].id] = v;
                  return acc;
                },
                {} as Record<string, Image>
              );
              setImages(() => rendered);
            }
          })
          .catch((err) => {
            console.log(err);
          });
      });
    }
  }, [searchParams]);

  const handleDelete = useCallback(() => {
    if (!formRef.current) {
      return;
    }
    const formData = new FormData(formRef.current);
    imagesToDelete.forEach((img) => {
      formData.append("delete_images[]", img);
      formData.append("delete_images_path[]", images[img].path);
    });
    console.log(formData);
    fetch("http://localhost:8888/api/similar", {
      method: "DELETE",
      body: formData,
    })
      .then((res) => {
        res.json().then((data: { deleted_images: string[] }) => {
          const newImages = Object.entries(images).filter(([k, v]) => {
            return !data.deleted_images.includes(k);
          });
          setImages(() => {
            return Object.fromEntries(newImages);
          });
        });
      })
      .catch((err) => {
        console.log(err);
      });
    clearCheckboxes();
  }, [imagesToDelete]);

  const handleSubmit = useCallback(
    (e: FormEvent) => {
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
          if (data.images) {
            const _images: Image[] = data.images;
            const rendered: Record<string, Image> = _images.reduce(
              (acc, v, idx) => {
                acc[_images[idx].id] = v;
                return acc;
              },
              {} as Record<string, Image>
            );
            console.log(rendered);
            setImages(() => rendered);
          }
          // console.log(data);
          // setImages(data.images || []);
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
    },
    [imagesToDelete, settings]
  );

  const handleMark = useCallback((id: string, checked: boolean) => {
    setImagesToDelete((prev) => {
      if (!checked) {
        return [...prev, id];
      } else {
        return prev.filter((uuid) => uuid !== id);
      }
    });
  }, []);

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
      <div tw="w-full">
        <form
          tw="items-center place-content-center flex flex-row ring-2 ring-slate-200 rounded-sm m-2 bg-[#333] "
          onSubmit={(e) => handleSubmit(e)}
          ref={formRef}
        >
          <InputGrid
            deleteHandler={handleDelete}
            applyHandler={handleApplySettings}
            settings={settings}
            setSettings={setSettings}
            imagesToDelete={imagesToDelete}
          />
          <div tw="flex flex-row gap-x-2"></div>
        </form>
      </div>
      <div tw="grid w-full h-auto mt-[5em] grid-flow-row  grid-cols-2 md:grid-cols-3 xl:grid-cols-6 items-end ">
        {Object.values(images).map((image, i) => {
          return (
            <MatchedImage
              key={i}
              src={`${image.path}`}
              alt={`similar-${i}`}
              image={image}
              uuid={image.id}
              checked={imagesToDelete.includes(image.id)}
              handleMark={handleMark}
              onClick={() => {
                setSearchParams({
                  uuid: image.id,
                  limit: settings.limit.toString(),
                  distance: settings.distance.toString(),
                });
              }}
            />
          );
        })}
      </div>
    </div>
  );
};

export default Similar;
