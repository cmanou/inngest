'use client';

import { Listbox } from '@headlessui/react';
import { NewButton } from '@inngest/components/Button';
import { OptionalTooltip } from '@inngest/components/Tooltip/OptionalTooltip';
import { RiArchive2Line, RiFirstAidKitLine, RiMore2Line } from '@remixicon/react';

export type AppActions = {
  isArchived: boolean;
  showArchive: () => void;
  showValidate: () => void;
  disableArchive?: boolean;
  disableValidate?: boolean;
};

export const ActionsMenu = ({
  isArchived,
  showArchive,
  showValidate,
  disableArchive = false,
  disableValidate = false,
}: AppActions) => {
  return (
    <Listbox>
      <Listbox.Button as="div">
        <NewButton kind="primary" appearance="outlined" size="medium" icon={<RiMore2Line />} />
      </Listbox.Button>
      <div className="relative">
        <Listbox.Options className="bg-canvasBase absolute right-1 top-5 z-50 w-[170px] gap-y-0.5 rounded border shadow">
          <Listbox.Option
            className="text-muted mx-2 mt-2 flex h-8 cursor-pointer items-center justify-start text-[13px]"
            value="eventKeys"
          >
            <OptionalTooltip
              tooltip={disableValidate && 'No syncs. App health check not available.'}
            >
              <NewButton
                disabled={disableValidate}
                appearance="ghost"
                kind="secondary"
                size="medium"
                icon={<RiFirstAidKitLine className="h-4 w-4" />}
                iconSide="left"
                label="Check app health"
                className={`text-muted m-0 w-full justify-start ${
                  disableValidate && 'cursor-not-allowed'
                }`}
                onClick={showValidate}
              />
            </OptionalTooltip>
          </Listbox.Option>

          {!isArchived && (
            <Listbox.Option
              className="m-2 flex h-8 cursor-pointer items-center text-[13px]"
              value="signingKeys"
            >
              <OptionalTooltip
                tooltip={disableArchive && 'Parent app is archived. Archive action not available.'}
              >
                <NewButton
                  appearance="ghost"
                  kind="danger"
                  size="medium"
                  icon={<RiArchive2Line className="h-4 w-4" />}
                  iconSide="left"
                  label={'Archive app'}
                  className="m-0 w-full justify-start"
                  onClick={showArchive}
                />
              </OptionalTooltip>
            </Listbox.Option>
          )}
        </Listbox.Options>
      </div>
    </Listbox>
  );
};
