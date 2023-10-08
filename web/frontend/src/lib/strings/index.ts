export const TrimURI = (uri: string) => {
  const lastSlashIndex = uri.lastIndexOf("/");
  const truncatedPath = uri.slice(0, lastSlashIndex);
  return truncatedPath;
};

export const roundedFloatFromString = (num: string, precision = 2) => {
  return (
    Math.round(parseFloat(num) * Math.pow(10, precision)) /
    Math.pow(10, precision)
  );
};


