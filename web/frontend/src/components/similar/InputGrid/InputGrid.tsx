import styled from "@emotion/styled";
import { roundedFloatFromString } from "lib/strings";
import tw from "twin.macro";

type Settings = {
  limit: number;
  distance: number;
  image_weight: number;
  text_weight: number;
  text_input: string;
  path: string;
  uuid: string;
};

export default function InputGrid({
  deleteHandler,
  applyHandler,
  setSettings,
  settings,
  imagesToDelete
}: {
  deleteHandler: () => void;
  applyHandler: () => void;
  setSettings: React.Dispatch<React.SetStateAction<Settings>>;
  settings: Settings;
  imagesToDelete: string[];
}) {
  return (
    <Container>
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
                    Math.round((1 - parseFloat(e.target.value)) * 100) / 100,
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
                    Math.round((1 - parseFloat(e.target.value)) * 100) / 100,
                };
              })
            }
          />
        </Field>
      </div>
      <div tw="flex flex-col gap-y-2">
        <Field>
          <label htmlFor="distance">Distance: {settings.distance}</label>
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
          <Button type="button" onClick={() => applyHandler()}>
            Apply
          </Button>
          <Button type="button" onClick={() => deleteHandler()}>
            Delete {imagesToDelete.length} images
          </Button>
        </div>
      </Field>
    </Container>
  );
}

const Container = styled("div")(() => [
  tw`grid grid-cols-[5fr 4fr 3fr 1fr ] m-2 p-2 w-full gap-x-2 place-content-around`,
]);

const Button = styled("button")(() => [
  tw`bg-slate-200 text-slate-800 rounded-sm p-1 m-1 cursor-pointer`,
]);

const SubmitButton = styled("input")(() => [
  tw`bg-slate-200 text-slate-800 rounded-sm p-1 m-1 cursor-pointer`,
]);

const Range = styled("input")(() => [tw`w-[7em]`]);

const Input = styled("input")(() => [tw`w-full`]);

const Field = styled("div")(() => [tw`flex flex-row justify-between`]);
