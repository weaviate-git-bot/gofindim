import {PrimaryButton} from "components/Buttons";
import { Fragment } from "react";
import { Link } from "react-router-dom";
import tw, { styled } from "twin.macro";

interface MatchedImageProps {
  src: string;
  alt: string;
  uuid: string;
  key: number;
  handleCheck: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onClick: (e: React.MouseEvent<HTMLImageElement>) => void;
}

const TrimURI = (uri: string) => {
    const lastSlashIndex = uri.lastIndexOf("/");
    const truncatedPath = uri.slice(0, lastSlashIndex);
    return truncatedPath;
}
const MatchedImage = ({
  src,
  alt,
  uuid,
  handleCheck,
  onClick,
}: MatchedImageProps) => {
  return (
    <Card>
      <Image tw="cursor-pointer" src={config.BaseURL+ src} alt={alt} onClick={onClick} />
      <div tw='flex flex-col w-fit items-center place-content-center '>
        <label htmlFor={uuid}>X</label>
        <div tw="flex flex-row">
          <input
            className="delete_images_check"
            type="checkbox"
            name="delete_images[]"
            value={uuid}
            onChange={handleCheck}
          />
        </div>

        <div tw='flex flex-row gap-x-2'>
        <Link target="_blank" to={`${config.BaseURL}/${src}`}><PrimaryButton>View</PrimaryButton></Link>
        <Link target="_blank"to={`/browse?path=${(TrimURI(src))}`}><PrimaryButton>Browse</PrimaryButton></Link>
        </div>
      </div>
    </Card>
  );
};

const Card = styled("div")(() => [
  tw`relative flex flex-col items-center justify-center  rounded-sm m-2`,
]);
const Image = styled("img")(() => [
  tw`w-full object-contain h-full max-w-[300px] max-h-[300px] flex items-center place-content-center`,
]);

export default MatchedImage;
