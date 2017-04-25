# LoggingSampleIntegrationServer

A very simple sample integration service for Testing/Logging:

* Enterprise Events
* Enterprise Validations
* Enterprise Numbering

## Events

Accepts document-style SOAP requests and logs to files in the "events" folder.

Sample Configuration:

```xml
<castle>

  <components>

    <component id="externalSoapEventChannel_{##UNIQUE##}" service="Eni.AfeNavigatorServer.Framework.Channels.IChannel, Eni.AfeNavigatorServer.Framework" type="Eni.AfeNavigatorServer.Framework.Channels.SimpleSoapChannel, Eni.AfeNavigatorServer.Framework"  lifestyle="singleton">
      <parameters>

        <serviceUrl>http://localhost:7878/event/process</serviceUrl>
        <serviceAction>processAfeEvent</serviceAction>
        <timeout>10</timeout>

        <!-- to use HTTP basic authentication, uncomment the following two lines and provider username/password -->
        <!-- recommended to only use this over SSL since basic auth doesn't encrypt passwords! -->
        <!-- <authUserName>user</authUserName> -->
        <!-- <authPassword>password</authPassword> -->

      </parameters>
    </component>

    <component id="externalSoapEventTransformer_{##UNIQUE##}" service="Eni.AfeNavigatorServer.Components.Events.IEventTransformer, Eni.AfeNavigatorServer.Components" type="Eni.AfeNavigatorServer.Components.Events.ClassicAfeEventTransformer, Eni.AfeNavigatorServer.Components"  lifestyle="transient">
    </component>

    <component id="historicalAfeEventFilter_{##UNIQUE##}" service="Eni.AfeNavigatorServer.Framework.Events.IEventFilter, Eni.AfeNavigatorServer.Framework" type="Eni.AfeNavigatorServer.Framework.Events.HistoricalAfeEventFilter, Eni.AfeNavigatorServer.Framework"  lifestyle="transient">
      <parameters>
        <suppressEventsOnHistoricalAFEs>false</suppressEventsOnHistoricalAFEs>
      </parameters>
    </component>

    <component id="eventFilter_{##UNIQUE##}" service="Eni.AfeNavigatorServer.Framework.Events.IEventFilter, Eni.AfeNavigatorServer.Framework" type="Eni.AfeNavigatorServer.Framework.Events.AggregatingEventFilter, Eni.AfeNavigatorServer.Framework"  lifestyle="transient">
      <parameters>
        <eventFilters>
          <array>
            <item>${historicalAfeEventFilter_{##UNIQUE##}}</item>
            <item>${eventTypeFilter_{##UNIQUE##}}</item>
          </array>
        </eventFilters>
      </parameters>
    </component>

    <component id="eventTypeFilter_{##UNIQUE##}" service="Eni.AfeNavigatorServer.Framework.Events.IEventFilter, Eni.AfeNavigatorServer.Framework" type="Eni.AfeNavigatorServer.Framework.Events.EventTypeFilter, Eni.AfeNavigatorServer.Framework"  lifestyle="transient">
      <parameters>
        <eventTypes>
          <array>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterInternalApprovalCompleteEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterInternalApprovalEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterInternalHoldEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterInternalRejectionEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterInternalUnHoldEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterFullApprovalEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterPartnerApprovalEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterPartnerRejectionEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterReleaseEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterReviewCompleteEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterReviewEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterRevisionEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterRouteForReviewEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterSupplementEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterRouteForReviewEvent</item>

<!--        <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterSaveEvent</item>
 -->
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.ManualExportEvent</item>

            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterSystemReviewEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterReviewerRemovedEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterReviewerAddedEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterAfeClosedEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterAfeReopenedEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterForceApprovalEvent</item>
            <item>Eni.AfeNavigatorServer.Framework.Events.AfeEvents.AfterUnreleaseEvent</item>
          </array>
        </eventTypes>
      </parameters>
    </component>

    <component id="externalEventHandler1_{##UNIQUE##}" service="Eni.AfeNavigatorServer.Framework.Events.IEventHandler, Eni.AfeNavigatorServer.Framework" type="Eni.AfeNavigatorServer.Components.Events.StandardEventHandler, Eni.AfeNavigatorServer.Components"  lifestyle="transient">
      <parameters>
        <channel>${externalSoapEventChannel_{##UNIQUE##}}</channel>
        <eventFilter>${eventFilter_{##UNIQUE##}}</eventFilter>
        <eventTransformer>${externalSoapEventTransformer_{##UNIQUE##}}</eventTransformer>
        <notifier>${smtpNotifier}</notifier>
        <emailAddresses>
          <array>
            <!-- List of email addresses to notify upon failures handling an event -->
          </array>
        </emailAddresses>
      </parameters>
    </component>

  </components>
</castle>
```

## Validations

Accepts validation requests and responds with canned response from `validation_reponse.xml` file.

Sample configuration:

```xml
<castle>
  <!-- ============================================================================== -->
  <!-- Expects to call a SOAP method (document style SOAP) that passes data
       conforming to afe-validate.xsd and returns a data conforming to afe-validate-result.xsd. -->
  <!-- ============================================================================== -->

  <components>
    <component id="validator_channel_{##UNIQUE##}" service="Eni.AfeNavigatorServer.Framework.Channels.IChannel, Eni.AfeNavigatorServer.Framework" type="Eni.AfeNavigatorServer.Framework.Channels.SimpleSoapChannel, Eni.AfeNavigatorServer.Framework" lifestyle="transient">
      <parameters>

        <serviceUrl>http://localhost:7878/validate/validate</serviceUrl>
        <serviceAction>validateAfe</serviceAction>
        <timeout>31</timeout>

      </parameters>
    </component>

    <component id="validator_transformer_{##UNIQUE##}" service="Eni.AfeNavigatorServer.Components.Validation.AfeValidation.IAfeValidationTransformer, Eni.AfeNavigatorServer.Components" type="Eni.AfeNavigatorServer.Components.Validation.AfeValidation.StandardAfeValidationTransformer, Eni.AfeNavigatorServer.Components"  lifestyle="transient"/>

    <component id="validator_{##UNIQUE##}" service="Eni.AfeNavigatorServer.Components.Validation.AfeValidation.ICustomAfeValidationRule, Eni.AfeNavigatorServer.Components" type="Eni.AfeNavigatorServer.Components.Validation.AfeValidation.CustomExternalAfeValidator, Eni.AfeNavigatorServer.Components"  lifestyle="transient">
      <parameters>
        <channel>${validator_channel_{##UNIQUE##}}</channel>
        <transformer>${validator_transformer_{##UNIQUE##}}</transformer>

        <!-- update these parameters to turn off validation for Route or Release -->
        <skipRoute>false</skipRoute>
        <skipRelease>false</skipRelease>

      </parameters>
    </component>
  </components>

</castle>
```

## Numbering

Accepts numbering requests and response with AFE Number based on current timestamp.

Sample configuration:

```xml
<castle>

  <components>
    <component id="afeNumberChannel" service="Eni.AfeNavigatorServer.Framework.Channels.IChannel, Eni.AfeNavigatorServer.Framework" type="Eni.AfeNavigatorServer.Framework.Channels.SimpleSoapChannel, Eni.AfeNavigatorServer.Framework"  lifestyle="singleton">
      <parameters>

        <serviceUrl>http://localhost:7878/numbering/getAfeNumber</serviceUrl>
        <serviceAction>getAfeNumber</serviceAction>
        <timeout>31</timeout>

      </parameters>
    </component>

    <component id="afeNumberRequestTransformer" service="Eni.AfeNavigatorServer.Components.AfeNumbering.IRequestTransformer, Eni.AfeNavigatorServer.Components" type="Eni.AfeNavigatorServer.Components.AfeNumbering.DefaultRequestTransformer, Eni.AfeNavigatorServer.Components"  lifestyle="transient" />

    <component id="afeNumberResponseParser" service="Eni.AfeNavigatorServer.Components.AfeNumbering.IResponseParser, Eni.AfeNavigatorServer.Components" type="Eni.AfeNavigatorServer.Components.AfeNumbering.DefaultResponseParser, Eni.AfeNavigatorServer.Components"  lifestyle="transient"/>

    <component id="afeNumberGenerator" service="Eni.AfeNavigatorServer.Components.AfeNumbering.IAfeNumberGenerator, Eni.AfeNavigatorServer.Components" type="Eni.AfeNavigatorServer.Components.AfeNumbering.ExternalAfeNumberGenerator, Eni.AfeNavigatorServer.Components"  lifestyle="transient">
      <parameters>
        <channel>${afeNumberChannel}</channel>
        <requestTransformer>${afeNumberRequestTransformer}</requestTransformer>
        <responseParser>${afeNumberResponseParser}</responseParser>
        <emailAddresses>
          <array>
            <!-- List of email addresses to notify upon failures handling an event -->
            <!-- <item>example@example.com</item> -->
          </array>
        </emailAddresses>
      </parameters>
    </component>

  </components>
</castle>
```