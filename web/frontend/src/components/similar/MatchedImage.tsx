import { PrimaryButton } from "@/components/Buttons";
import { TrimURI } from "@/lib/strings";
import { Link, useSearchParams } from "react-router-dom";
import { Image as ImageType } from "@/types/Image";
import tw, { styled } from "twin.macro";

interface MatchedImageProps {
  src: string;
  alt: string;
  uuid: string;
  image: ImageType;
  checked: boolean;
  handleMark: (id: string, checked: boolean) => void;
  onClick: (e: React.MouseEvent<HTMLImageElement>) => void;
}

const MatchedImage = ({
  alt,
  checked,
  image,
  handleMark,
  onClick,
}: MatchedImageProps) => {
  const [searchParams, setSearchParams] = useSearchParams();
  return (
    <Card checked={checked}>
    <span tw='absolute top-0 left-[50%] translate-x-[-50%] translate-y-[-100%]'>
    {searchParams.get("uuid")==image.id && "Active"}
    </span>
      <Image
        tw="cursor-pointer"
        src={config.BaseURL + image.path}
        alt={alt}
        onClick={onClick}
      />
      <input
        hidden
        className="delete_images_check"
        type="checkbox"
        name="delete_images[]"
        value={image.id}
        checked={checked}
      />
      <div tw="flex flex-col w-fit items-center place-content-center ">
        <div tw="flex flex-row gap-x-2">
          <Link target="_blank" to={`${config.BaseURL}/${image.path}`}>
            <PrimaryButton>View</PrimaryButton>
          </Link>
          <Link target="_blank" to={`/browse?path=${TrimURI(image.path)}`}>
            <PrimaryButton>Browse</PrimaryButton>
          </Link>
          <PrimaryButton
            onClick={() => handleMark(image.id, checked)}
            tw="bg-red-700"
          >
            {checked ? "Unmark" : "Mark"}
          </PrimaryButton>
        </div>
      </div>
    </Card>
  );
};

const Card = styled("div")(({ checked }: { checked: boolean }) => [
  tw`relative flex flex-col items-center justify-center rounded-sm m-2 h-full`,
  checked && tw`ring-2 ring-slate-200`,
]);
const Image = styled("img")(() => [
  tw`w-full object-contain h-full max-w-[300px] max-h-[300px] flex items-center place-content-center`,
]);

export default MatchedImage;
