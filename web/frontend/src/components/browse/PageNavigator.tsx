import ReactPaginate from "react-paginate";
import tw, {css} from "twin.macro";

interface Props {
  page: number;
  setPage: (page: number) => void;
  pageCount: number;
  pageSize: number;
}

const PageNavigator = ({ page, setPage, pageCount, pageSize }: Props) => {
  const handlePageClick = (event: { selected: number }) => {
    const selectedPage = event.selected;
    setPage(selectedPage);
  };
  return (
    <div tw=' top-10 flex items-center place-content-around gap-x-2 flex-row justify-around'>
      <ReactPaginate
      tw='flex flex-row w-full justify-between items-center place-content-center'
      css={css`
      .active-page {
        ${tw`bg-slate-500 px-2 py-1 rounded-full`}
      }
      `}
      activeClassName="active-page"
        breakLabel="..."
        nextLabel="next >"
        onPageChange={handlePageClick}
        pageRangeDisplayed={5}
        pageCount={pageCount}
        previousLabel="< previous"
        renderOnZeroPageCount={null}
      />
    </div>
  );
};

export default PageNavigator; 
